async function postJSON(url, data) {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  const json = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(json.error || `HTTP ${res.status}`);
  }
  return json;
}

function setMsg(text) {
  const el = document.getElementById("msg");
  if (el) el.textContent = text;
}

document.addEventListener("DOMContentLoaded", () => {
  const loginForm = document.getElementById("loginForm");
  const registerForm = document.getElementById("registerForm");
  const forgotForm = document.getElementById("forgotForm");
  const resetForm = document.getElementById("resetForm");  

  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const fd = new FormData(loginForm);
      try {
        const resp = await postJSON("/api/login", {
          email: fd.get("email"),
          password: fd.get("password"),
        });
        // resp.data.token from API
        localStorage.setItem("token", resp.data.token);
        localStorage.setItem("email", resp.data.user);
        window.location.href = "/profile";
      } catch (err) {
        setMsg(String(err.message || err));
      }
    });
  }

  if (registerForm) {
    registerForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const fd = new FormData(registerForm);
      try {
        const resp = await postJSON("/api/register", {
          email: fd.get("email"),
          password: fd.get("password"),
        });
        localStorage.setItem("token", resp.data.token);
        localStorage.setItem("email", resp.data.user);
        window.location.href = "/profile";
      } catch (err) {
        setMsg(String(err.message || err));
      }
    });
  }

  if (forgotForm) {
    forgotForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      setMsg("");

      const fd = new FormData(forgotForm);
      const email = fd.get("email");

      try {
        const resp = await postJSON("/api/forgot-password", { email });
        // API response
        const token = resp?.data?.token;
        if (!token) throw new Error("No token returned from API");

        // Redirect with prefilled query
        window.location.href = `/reset-password?token=${encodeURIComponent(token)}`;
      } catch (err) {
        setMsg(String(err.message || err));
      }
    });
  }

  if (resetForm) {
    resetForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      setMsg("");

      const fd = new FormData(resetForm);
      const token = fd.get("token");
      const password = fd.get("password");

      try {
        await postJSON("/api/reset-password", { token, password });
        setMsg("Password reset OK. Redirecting to login...");
        setTimeout(() => {
          window.location.href = "/login";
        }, 800);
      } catch (err) {
        setMsg(String(err.message || err));
      }
    });
  }
});