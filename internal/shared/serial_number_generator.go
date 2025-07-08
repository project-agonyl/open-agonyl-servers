package shared

import (
	"context"
	"fmt"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type SerialNumberGenerator interface {
	GetNextSerial(ctx context.Context) (uint32, error)
	GetNextSerials(ctx context.Context, count int) ([]uint32, error)
	GetSequenceInfo(ctx context.Context) (*ItemSequence, error)
	GetAllSequences(ctx context.Context) ([]ItemSequence, error)
	CleanupOldSequences(ctx context.Context, olderThan time.Time) error
}

const batchSize = 500

type serialNumberGenerator struct {
	db         *sqlx.DB
	redis      *redis.Client
	serverID   string
	batchSize  int32
	mu         sync.Mutex
	currentID  uint32
	maxID      uint32
	lockKey    string
	counterKey string
	psql       sq.StatementBuilderType
}

func NewSerialNumberGenerator(db *sqlx.DB, redisClient *redis.Client, serverID string) SerialNumberGenerator {
	return &serialNumberGenerator{
		db:         db,
		redis:      redisClient,
		serverID:   serverID,
		batchSize:  batchSize,
		lockKey:    fmt.Sprintf("serial:lock:%s", serverID),
		counterKey: fmt.Sprintf("serial:counter:%s", serverID),
		psql:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (sng *serialNumberGenerator) GetNextSerial(ctx context.Context) (uint32, error) {
	sng.mu.Lock()
	defer sng.mu.Unlock()
	if sng.currentID >= sng.maxID {
		if err := sng.allocateNewBatch(ctx); err != nil {
			return 0, err
		}
	}

	sng.currentID++
	return sng.currentID, nil
}

func (sng *serialNumberGenerator) GetNextSerials(ctx context.Context, count int) ([]uint32, error) {
	sng.mu.Lock()
	defer sng.mu.Unlock()
	serials := make([]uint32, 0, count)
	for i := 0; i < count; i++ {
		if sng.currentID >= sng.maxID {
			if err := sng.allocateNewBatch(ctx); err != nil {
				return serials, err
			}
		}

		sng.currentID++
		serials = append(serials, sng.currentID)
	}

	return serials, nil
}

func (sng *serialNumberGenerator) GetSequenceInfo(ctx context.Context) (*ItemSequence, error) {
	query := sng.psql.Select("*").
		From("item_sequences").
		Where(sq.Eq{"server_id": sng.serverID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var sequence ItemSequence
	err = sng.db.GetContext(ctx, &sequence, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence info: %w", err)
	}

	return &sequence, nil
}

func (sng *serialNumberGenerator) GetAllSequences(ctx context.Context) ([]ItemSequence, error) {
	query := sng.psql.Select("*").
		From("item_sequences").
		OrderBy("server_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var sequences []ItemSequence
	err = sng.db.SelectContext(ctx, &sequences, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequences: %w", err)
	}

	return sequences, nil
}

func (sng *serialNumberGenerator) CleanupOldSequences(ctx context.Context, olderThan time.Time) error {
	query := sng.psql.Delete("item_sequences").
		Where(sq.Lt{"updated_at": olderThan})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build cleanup query: %w", err)
	}

	result, err := sng.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to cleanup old sequences: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Cleaned up %d old sequence records\n", rowsAffected)

	return nil
}

func (sng *serialNumberGenerator) allocateNewBatch(ctx context.Context) error {
	lockDuration := 5 * time.Second
	acquired, err := sng.redis.SetNX(ctx, sng.lockKey, sng.serverID, lockDuration).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		return fmt.Errorf("failed to acquire lock: already held")
	}

	defer sng.redis.Del(ctx, sng.lockKey)
	if start, end, err := sng.getCachedBatch(ctx); err == nil {
		sng.currentID = start - 1
		sng.maxID = end
		return nil
	}

	batch, err := sng.allocateBatchFromDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to allocate batch: %w", err)
	}

	pipe := sng.redis.Pipeline()
	pipe.HSet(ctx, sng.counterKey, "start", batch.StartID, "end", batch.EndID)
	pipe.Expire(ctx, sng.counterKey, 10*time.Minute)
	_, err = pipe.Exec(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to cache batch: %v\n", err)
	}

	sng.currentID = uint32(batch.StartID - 1)
	sng.maxID = uint32(batch.EndID)
	return nil
}

func (sng *serialNumberGenerator) allocateBatchFromDB(ctx context.Context) (*BatchAllocation, error) {
	query := `SELECT start_id, end_id FROM allocate_sequence_batch($1, $2)`
	var batch BatchAllocation
	err := sng.db.GetContext(ctx, &batch, query, sng.serverID, sng.batchSize)
	if err != nil {
		return nil, err
	}

	return &batch, nil
}

func (sng *serialNumberGenerator) getCachedBatch(ctx context.Context) (uint32, uint32, error) {
	result := sng.redis.HMGet(ctx, sng.counterKey, "start", "end")
	if result.Err() != nil {
		return 0, 0, result.Err()
	}

	vals := result.Val()
	if len(vals) != 2 || vals[0] == nil || vals[1] == nil {
		return 0, 0, fmt.Errorf("incomplete cache data")
	}

	startStr, ok1 := vals[0].(string)
	endStr, ok2 := vals[1].(string)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("invalid cache data types")
	}

	var start, end int64
	if _, err := fmt.Sscanf(startStr, "%d", &start); err != nil {
		return 0, 0, err
	}

	if _, err := fmt.Sscanf(endStr, "%d", &end); err != nil {
		return 0, 0, err
	}

	return uint32(start), uint32(end), nil
}

type ItemSequence struct {
	ID            int       `db:"id"`
	ServerID      string    `db:"server_id"`
	LastAllocated int64     `db:"last_allocated"`
	BatchSize     int32     `db:"batch_size"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type BatchAllocation struct {
	StartID int64 `db:"start_id"`
	EndID   int64 `db:"end_id"`
}
