async function fetchCurrentUser() {
    const res = await fetch("/api/me", { credentials: "include" });
    if (!res.ok) return null;
    return await res.json();
}

document.addEventListener("DOMContentLoaded", async () => {
    const user = await fetchCurrentUser();
    if (user) {
        document.querySelector("#welcome").textContent = `Welcome, ${user.username}!`;
    } else {
        window.location.href = "/login.html";
    }
});