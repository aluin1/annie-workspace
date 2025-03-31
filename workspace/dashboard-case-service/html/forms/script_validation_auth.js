import CONFIG from './config_data_case.js';
 
const URL_VALIDATION_TOKEN_GMAIL =  CONFIG.URL_VALIDATION_TOKEN_GMAIL;


 

export async function handleCredentialResponse(response) {
    //console.log("📥 Google JWT Token:", response.credential);
    await ValidationToken(response.credential);
}

 
async function ValidationToken(tokenGmail) {

    // Siapkan request untuk dikirim ke backend
    const requestData = { token: tokenGmail };
    console.log("📤 Request Validation Auth Token:", requestData);

    try {
        const response = await fetch(`${URL_VALIDATION_TOKEN_GMAIL}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(requestData),
        });

        
        const data = await response.json();

        if (!response.ok) {
            throw new Error(`${response.status} - ${response.statusText} \n ${data.response_message}`);
        }


        if (data && data.email) {
            console.log("✅ Token Valid Account:", data.email); 
            localStorage.setItem("google_jwt_token", tokenGmail);
            localStorage.setItem("email", data.email);
            localStorage.setItem("name", data.name);
            localStorage.setItem("picture", data.picture);

            await swal({
                title: "Success",
                text: `Welcome, ${data.name}!`,
                icon: "success",
                button: "OK"
            });

            window.location.href = "data_case.html"; // Redirect setelah klik OK
            return true; // Token valid
        } else {
            throw new Error(data.response_message);
        }
    } catch (err) {
        console.error("❌ Error validate token:", err);
        // swal.close();

        await swal({
            title: "Credential Failed",
            text: err.message,
            icon: "error",
            buttons: { confirm: { className: "btn btn-danger" } }
        });
 
        window.location.href = "login.html";
        return false; // Token tidak valid
    }
}
