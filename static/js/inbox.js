function toggleInbox() {
    const dropdown = document.getElementById('inboxDropdown');
    if (dropdown.style.display === 'none' || dropdown.style.display === '') {
        dropdown.style.display = 'block';
        loadInbox();
    } else {
        dropdown.style.display = 'none';
    }
}

async function loadInbox() {
    const content = document.getElementById('inboxContent');
    
    try {
        const response = await fetch('/api/friend-request/pending', {
            credentials: 'include'
        });
        
        if (!response.ok) {
            content.textContent = 'Failed to load notifications';
            return;
        }
        
        const requests = await response.json() || [];
        
        if (requests.length === 0) {
            content.textContent = 'You have no notifications';
            return;
        }
        
        // Build HTML for each request
        let html = '';
        console.log('Requests received:', requests); // Add this before the forEach
        requests.forEach(req => {
            console.log('Request object:', req);
            html += `
                <div style="display: flex; align-items: center; justify-content: space-between; padding: 5px 0; border-bottom: 1px solid #eee;">
                    <span>${req.Username}: Added you as friend.</span>
                    <div style="display: flex; gap: 5px;">
                        <button onclick="handleRequest('${req.FriendshipID }', false)" style="padding: 2px 6px; cursor: pointer; background: none; border: none;">
                            <img src="/static/assets/reject.png" alt="reject" style="width: 16px; height: 16px;">
                        </button>
                        <button onclick="handleRequest('${req.FriendshipID }', true)" style="padding: 2px 6px; cursor: pointer; background: none; border: none;">
                            <img src="/static/assets/accept.png" alt="accept" style="width: 16px; height: 16px;">
                        </button>
                    </div>
                </div>
            `;
        });
        
        content.innerHTML = html;
        
    } catch (error) {
        console.error('Error loading inbox:', error);
        content.textContent = 'Error loading notifications';
    }
}

function handleRequest(friendshipId, accept) {
    const status = accept ? 'accepted' : 'rejected';
    const action = accept ? 'Accepting' : 'Rejecting';
    
    console.log(`${action} friend request:`, friendshipId);
    
    fetch(`/api/friend-request`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
            friendship_id: friendshipId,
            status: status
        })
    })
    .then(res => {
        if (res.ok) {
            alert(`Friend request ${accept ? 'accepted' : 'rejected'}!`);
            // Refresh the inbox to remove the request
            loadInbox();
            if (accept && typeof loadFriends === 'function') {
                loadFriends();
            }
        } else {
            return res.json().then(data => {
                alert(data.error || `Failed to ${accept ? 'accept' : 'reject'} friend request`);
            });
        }
    })
    .catch(err => {
        console.error('Error:', err);
        alert('Error processing friend request');
    });
}

// Close dropdown when clicking outside
document.addEventListener('click', function(event) {
    const inboxContainer = document.querySelector('.inbox-container');
    if (inboxContainer && !inboxContainer.contains(event.target)) {
        const dropdown = document.getElementById('inboxDropdown');
        if (dropdown) {
            dropdown.style.display = 'none';
        }
    }
});