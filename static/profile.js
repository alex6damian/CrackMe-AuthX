document.addEventListener("DOMContentLoaded", async () => {
  const out = document.getElementById("out");
  const token = localStorage.getItem("token");

  if (!token) {
    window.location.href = "/login";
    return;
  }

  const res = await fetch("/api/profile", {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (res.status === 401) {
    window.location.href = "/login";
    return;
  }

  const json = await res.json();
  out.textContent = JSON.stringify(json, null, 2);

  document.getElementById("logoutBtn").addEventListener("click", () => {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    localStorage.removeItem("created_at");
    window.location.href = "/login";
  });
});