import CONFIG from './config_data_case.js';

const URL_VALIDATION_TOKEN_GMAIL = CONFIG.URL_VALIDATION_TOKEN_GMAIL;

// Fungsi utama dari Google Sign-In
export async function handleCredentialResponse(response) {
    if (response && response.credential) {
        await validateToken(response.credential);
    } else {
        console.error("❌ No credential received.");
    }
}

// Validasi token ke backend
async function validateToken(tokenGmail) {
    const requestData = { token: tokenGmail };
    console.log("📤 Sending token for validation:", requestData);

    swal({ title: "Validating token...", text: "Please wait...", icon: "info", buttons: false });

    try {
        const response = await fetch(URL_VALIDATION_TOKEN_GMAIL, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(requestData),
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(`${response.status} ${response.statusText}\n${data.response_message}`);
        }

        if (data && data.email) {
            console.log("✅ Token valid:", data.email);

            localStorage.setItem("google_jwt_token", tokenGmail);
            localStorage.setItem("email", data.email);
            localStorage.setItem("name", data.name);
            localStorage.setItem("picture", data.picture);

            swal.close();

            await swal({
                title: "Success",
                text: `Welcome, ${data.name}!`,
                icon: "success",
                button: "OK"
            });

            window.location.href = "data_case.html";
            return true;
        } else {
            throw new Error(data.response_message || "Unknown server response.");
        }

    } catch (err) {
        console.error("❌ Token validation failed:", err.message);
        swal.close();

        await swal({
            title: "Login Failed",
            text: err.message,
            icon: "error",
            button: "OK"
        });

        window.location.href = "login.html";
        return false;
    }
}
