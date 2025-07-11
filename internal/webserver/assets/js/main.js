document.addEventListener('DOMContentLoaded', function() {
    console.log('A3 Agonyl Web Server loaded');
    
    initializeApp();
});

function initializeApp() {
    setupHTMXEvents();
    setupNavigationHighlight();
}

function setupHTMXEvents() {
    document.body.addEventListener('htmx:beforeRequest', function(event) {
        console.log('HTMX request starting:', event.detail.xhr.responseURL);
        showLoadingIndicator();
    });
    
    document.body.addEventListener('htmx:afterRequest', function(event) {
        console.log('HTMX request completed:', event.detail.xhr.responseURL);
        hideLoadingIndicator();
    });
    
    document.body.addEventListener('htmx:responseError', function(event) {
        console.error('HTMX request error:', event.detail.xhr.responseURL);
        showErrorMessage('Request failed. Please try again.');
        hideLoadingIndicator();
    });
}

function setupNavigationHighlight() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('nav a');
    
    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        }
    });
}

function showLoadingIndicator() {
    const indicator = document.querySelector('.htmx-indicator');
    if (indicator) {
        indicator.style.opacity = '1';
    }
}

function hideLoadingIndicator() {
    const indicator = document.querySelector('.htmx-indicator');
    if (indicator) {
        indicator.style.opacity = '0';
    }
}

function showErrorMessage(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'alert alert-danger';
    errorDiv.textContent = message;
    errorDiv.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 1000;
        padding: 15px;
        background-color: #f8d7da;
        border: 1px solid #f5c6cb;
        border-radius: 4px;
        color: #721c24;
    `;
    
    document.body.appendChild(errorDiv);
    
    setTimeout(() => {
        if (errorDiv.parentNode) {
            errorDiv.parentNode.removeChild(errorDiv);
        }
    }, 5000);
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

window.A3Agonyl = {
    showErrorMessage,
    debounce
}; 