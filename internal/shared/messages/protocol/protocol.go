package protocol

const C2SLogin uint16 = 0xE0
const C2SServerDetails uint16 = 0xE1
const S2CLoginMessage uint16 = 0xE2
const S2CServerList uint16 = 0xE3
const S2CLoginOk uint16 = 0xE4
const S2CServerDetails uint16 = 0xE5

const S2CError uint16 = 0x0FFF
const C2SKeepAlive uint16 = 0x0FF2
const C2SUnknownProtocol uint16 = 0x000F

const S2CPcAppear uint16 = 0x1100
const S2CPcDisappear uint16 = 0x1101
const C2SPrepareUser uint16 = 0x1105
const S2CCharacterList uint16 = 0x1105
const C2SCharacterLogin uint16 = 0x1106
const S2CCharacterLoginOk uint16 = 0x1106
const C2SWorldLogin uint16 = 0x1107
const S2CWorldLogin uint16 = 0x1107
const C2SCharacterLogout uint16 = 0x1108
const S2CCharLogout uint16 = 0x1108
const C2SWarp uint16 = 0x1111
const C2SReturn2Here uint16 = 0x1112
const C2SSubmapInfo uint16 = 0x1114
const C2SEnter uint16 = 0x1115
const C2SActivePet uint16 = 0x11A1
const C2SInactivePet uint16 = 0x11A2
const C2SPetBuy uint16 = 0x11A5
const C2SPetSell uint16 = 0x11A6
const C2SFeedPet uint16 = 0x11A7
const C2SRevivePet uint16 = 0x11A8
const C2SShueCombination uint16 = 0x11B0

const C2SAskMove uint16 = 0x1200
const S2CAnsMove uint16 = 0x1200
const S2CSeeMove uint16 = 0x1201
const C2SPcMove uint16 = 0x1202
const S2CFixMove uint16 = 0x1203
const S2CSeeStop uint16 = 0x1204
const C2SAskHsMove uint16 = 0x1205
const C2SHsMove uint16 = 0x1208

const S2CNpcInitializeProtocol uint16 = 0x1300
const C2SObjectNpc uint16 = 0x1307
const C2SAskNpcFavor uint16 = 0x1308
const C2SNpcFavorUp uint16 = 0x1309

const C2SAskAttack uint16 = 0x1400
const C2SLearnSkill uint16 = 0x1451
const C2SAskSkill uint16 = 0x1453
const C2SSkillSlotInfo uint16 = 0x1461
const S2CSkillSlotInfo uint16 = 0x1461
const C2SAnsRecall uint16 = 0x1462

const C2SAllotPoint uint16 = 0x1602
const C2SAskHeal uint16 = 0x1606
const C2SRetrievePoint uint16 = 0x1609
const C2SRestoreExp uint16 = 0x160C
const S2CUnknown37Protocol uint16 = 0x1610
const C2SLearnPskill uint16 = 0x1611
const C2SForgetAllPskill uint16 = 0x1613
const C2SAskOpenStorage uint16 = 0x1651
const C2SAskInven2Storage uint16 = 0x1652
const C2SAskStorage2Inven uint16 = 0x1653
const C2SAskDepositeMoney uint16 = 0x1654
const C2SAskWithdrawMoney uint16 = 0x1655
const C2SAskCloseStorage uint16 = 0x1656
const C2SAskMoveItemInStorage uint16 = 0x1657

const C2SPickupItem uint16 = 0x1702
const C2SDropItem uint16 = 0x1704
const C2SMoveItem uint16 = 0x1706
const C2SWearItem uint16 = 0x1708
const C2SStripItem uint16 = 0x1711
const C2SBuyItem uint16 = 0x1714
const C2SSellItem uint16 = 0x1716
const C2SGiveItem uint16 = 0x1718
const S2CGiveItem uint16 = 0x1720
const C2SUsePotion uint16 = 0x1721
const C2SAskDeal uint16 = 0x1723
const S2CAskDeal uint16 = 0x1724
const C2SAnsDeal uint16 = 0x1725
const C2SPutInItem uint16 = 0x1727
const S2CPutInItem uint16 = 0x1728
const C2SPutOutItem uint16 = 0x1729
const S2CPutOutItem uint16 = 0x1730
const C2SDecideDeal uint16 = 0x1731
const C2SConfirmDeal uint16 = 0x1733
const C2SUseItem uint16 = 0x1736
const C2SConfirmItem uint16 = 0x1742
const C2SRemodelItem uint16 = 0x1744
const C2SUseScroll uint16 = 0x1748
const C2SPutInPet uint16 = 0x1750
const C2SPutOutPet uint16 = 0x1751
const C2SItemCombination uint16 = 0x1753
const C2SLottoPurchase uint16 = 0x1754
const C2SLottoQueryPrize uint16 = 0x1755
const C2SLottoQueryHistory uint16 = 0x1756
const C2SLottoSale uint16 = 0x1757
const C2STakeItemInBox uint16 = 0x1760
const C2STakeItemOutBox uint16 = 0x1761
const C2SUsePotionEx uint16 = 0x1767
const C2SOpenMarket uint16 = 0x1770
const C2SCloseMarket uint16 = 0x1771
const C2SEnterMarket uint16 = 0x1773
const C2SBuyItemMarket uint16 = 0x1775
const C2SLeaveMarket uint16 = 0x1776
const C2SModifyMarket uint16 = 0x1777
const C2SAskItemSerial uint16 = 0x1780
const C2SSocketItem uint16 = 0x1781
const C2SBuyBattlefieldItem uint16 = 0x1785
const C2SBuyCashItem uint16 = 0x1790
const C2SCashInfo uint16 = 0x1791
const C2SDerbyIndexQuery uint16 = 0x17A9
const C2SDerbyMonsterQuery uint16 = 0x17AA
const C2SDerbyRatioQuery uint16 = 0x17AB
const C2SDerbyPurchase uint16 = 0x17AC
const C2SDerbyResultQuery uint16 = 0x17AE
const C2SDerbyHistoryQuery uint16 = 0x17AF
const C2SDerbyExchange uint16 = 0x17B0

const C2SSay uint16 = 0x1800
const C2SGesture uint16 = 0x1801
const C2SChatWindowOpt uint16 = 0x1803
const S2CChatWindowOpt uint16 = 0x1803

const C2SOption uint16 = 0x1900

const C2SPartyQuest uint16 = 0x2110
const C2SQuestExDialogueReq uint16 = 0x2120
const C2SQuestExDialogueAns uint16 = 0x2122
const C2SQuestExCancel uint16 = 0x2126
const C2SQuestExList uint16 = 0x2128
const C2SSquestStart uint16 = 0x2140
const C2SSquestStepEnd uint16 = 0x2144
const C2SSquestHistory uint16 = 0x2145
const C2SSquestMinigameMove uint16 = 0x2149
const C2SSquestWallQuiz uint16 = 0x214A
const C2SSquestWallOk uint16 = 0x214C
const C2SSquestA3QuizSelect uint16 = 0x214E
const C2SSquestA3Quiz uint16 = 0x214F
const C2SSquestA3QuizOk uint16 = 0x2151
const C2SSquestEndOk uint16 = 0x2152
const C2SSquest222NumQuiz uint16 = 0x2153
const C2SSquest312ItemCreate uint16 = 0x2155
const C2SSquestHboyRune uint16 = 0x2157
const C2SSquestHboyHanoi uint16 = 0x2159
const C2SSquest346ItemCombi uint16 = 0x215E

const C2SAskParty uint16 = 0x2200
const C2SAnsParty uint16 = 0x2202
const C2SOutParty uint16 = 0x2205
const C2SAskApprenticeIn uint16 = 0x22A0
const C2SAnsApprenticeIn uint16 = 0x22A1
const C2SAskApprenticeOut uint16 = 0x22A4

const C2SClan uint16 = 0x2300
const C2SJoinClan uint16 = 0x2301
const C2SAnsClan uint16 = 0x2302
const C2SBoltClan uint16 = 0x2303
const C2SReqClanInfo uint16 = 0x2304
const C2ZRegisterMark uint16 = 0x2320
const C2STransferMark uint16 = 0x2322
const C2SAskMark uint16 = 0x2323
const C2SFriendInfo uint16 = 0x2331
const C2SFriendState uint16 = 0x2332
const S2CFriendState uint16 = 0x2332
const C2SFriendGroup uint16 = 0x2333
const C2SAskFriend uint16 = 0x2334
const C2SAnsFriend uint16 = 0x2335
const C2SAskClanBattle uint16 = 0x2340
const C2SAnsClanBattle uint16 = 0x2341
const C2SAskClanBattleEnd uint16 = 0x2342
const C2SAnsClanBattleEnd uint16 = 0x2343
const C2SAskClanBattleScore uint16 = 0x2345
const C2SLetterBaseInfo uint16 = 0x2350
const C2SLetterSimpleInfo uint16 = 0x2351
const C2SLetterDel uint16 = 0x2353
const C2SLetterSend uint16 = 0x2354
const C2SLetterKeeping uint16 = 0x2356

const C2SChangeNation uint16 = 0x2400

const C2SCaoMitigation uint16 = 0x2510

const C2SAgitInfo uint16 = 0x2600
const C2SAuctionInfo uint16 = 0x2601
const C2SAgitEnter uint16 = 0x2602
const C2SAgitPutUpAuction uint16 = 0x2603
const C2SAgitBidOn uint16 = 0x2604
const C2SAgitPayExpense uint16 = 0x2605
const C2SAgitChangeName uint16 = 0x2606
const C2SAgitRepayMoney uint16 = 0x2607
const C2SAgitObtainSaleMoney uint16 = 0x2608
const C2SAgitManageInfo uint16 = 0x260A
const C2SAgitOption uint16 = 0x260B
const C2SAgitOptionInfo uint16 = 0x260C
const C2SAgitPcBan uint16 = 0x260D

const C2SChristmasCard uint16 = 0x2730
const C2SSpeakCard uint16 = 0x2731
const C2SProcessInfo uint16 = 0x2740

const C2SAskWarpZ2B uint16 = 0x3500
const C2SAskWarpB2Z uint16 = 0x3510

const C2SAskShopInfo uint16 = 0x3915
const C2SAskGiveMyTax uint16 = 0x3916

const C2STyrUnitList uint16 = 0x4001
const C2STyrUnitInfo uint16 = 0x4002
const C2STyrEntry uint16 = 0x4003
const C2STyrJoin uint16 = 0x4004
const C2STyrRewardInfo uint16 = 0x4080
const C2STyrReward uint16 = 0x4081

const C2STyrUpgrade uint16 = 0x4102

const C2STyrRtmmEnd uint16 = 0x4203

const C2SHsSeal uint16 = 0x5001
const C2SHsRecall uint16 = 0x5002
const C2SHsRevive uint16 = 0x5005
const C2SHsAskAttack uint16 = 0x5006
const C2SHsStoneBuy uint16 = 0x5008
const C2SHsStoneSell uint16 = 0x5009
const C2SHsLearnSkill uint16 = 0x500A
const C2SHsAllotPoint uint16 = 0x500B
const C2SHsRetrievePoint uint16 = 0x500C
const C2SHsWearItem uint16 = 0x500D
const C2SHsStripItem uint16 = 0x5010
const C2SHsOption uint16 = 0x501B
const C2SHsHeal uint16 = 0x501C
const C2SHsSkillReset uint16 = 0x501E

const C2SAskMigration uint16 = 0x9000

const M2SError uint16 = 0xA000
const C2SAskCreatePlayer uint16 = 0xA001
const S2CAnsCreatePlayer uint16 = 0xA001
const C2SAskDeletePlayer uint16 = 0xA002
const S2CAnsDeletePlayer uint16 = 0xA002
const S2MCharacterLogin uint16 = 0xA010

const C2SLeague uint16 = 0xA340
const C2SReqLeagueClanInfo uint16 = 0xA345
const C2SLeagueAllow uint16 = 0xA347

const C2SPayInfo uint16 = 0xC000
const S2MMapList uint16 = 0xC001

const C2SPing uint16 = 0xF001
