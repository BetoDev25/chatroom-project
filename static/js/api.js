let currentUser = null;

async function fetchCurrentUser() {
    const res = await fetch("/api/me", { credentials: "include" });
    if (!res.ok) return null;
    currentUser = await res.json();
    return currentUser;
}

document.addEventListener("DOMContentLoaded", async () => {
    const user = await fetchCurrentUser();
    if (user) {
        document.querySelector("#welcome").textContent = `Welcome, ${user.username}!`;
        window.currentUser = user;
    } else {
        window.location.href = "/login.html";
    }
});
