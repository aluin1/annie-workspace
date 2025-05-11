// const CONFIG = {
//     URL_INSERT_CASE: 'http://localhost:8081/insertCase',
//     URL_GET_CASE: 'http://localhost:8081/dataCase',
//     URL_GET_TOKEN: 'http://localhost:8081/getToken', 
//     URL_EDIT_STATUS_CASE: 'http://localhost:8081/editCaseStatus', 
//     CLIENT_ID:'79d4270a98d241999fa7e54e015b51ab',
//     CLIENT_SECRET:'22106803111c4cdda620875ed4779129' ,
//     // URL_UPLOAD_FILE:'http://localhost:8081/uploadFiles',
//     URL_UPLOAD_FILE:'http://localhost:8081/UploadFilesCloud',
//     EMAIL_SUPPORT:'fmbods@bioprogressive.org',
//     URL_VALIDATION_TOKEN_GMAIL:'http://localhost:8081/validationTokenGmail',
//     URL_GET_DATA_ANNIE: 'http://localhost:8081/GetDataAnnie' 
// };

const CONFIG = {
    URL_INSERT_CASE: 'https://annievipbe.ismartdds.com/insertCase',
    URL_GET_CASE: 'https://annievipbe.ismartdds.com/dataCase',
    URL_GET_TOKEN: 'https://annievipbe.ismartdds.com/getToken', 
    URL_EDIT_STATUS_CASE: 'https://annievipbe.ismartdds.com/editCaseStatus', 
    CLIENT_ID:'79d4270a98d241999fa7e54e015b51ab',
    CLIENT_SECRET:'22106803111c4cdda620875ed4779129' ,
    // URL_UPLOAD_FILE:'https://annievipbe.ismartdds.com/uploadFiles',
    URL_UPLOAD_FILE:'https://annievipbe.ismartdds.com/UploadFilesCloud',
    EMAIL_SUPPORT:'fmbods@bioprogressive.org',
    URL_VALIDATION_TOKEN_GMAIL:'https://annievipbe.ismartdds.com/validationTokenGmail',
    URL_GET_DATA_ANNIE: 'https://annievipbe.ismartdds.com/GetDataAnnie'
};

document.addEventListener("DOMContentLoaded", checkDataUser);

async function checkDataUser() {
    console.log("🚀 DOM Loaded - Checking authentication...");

    const userEmail = localStorage.getItem("email") || "";
    const userPicture = localStorage.getItem("picture") || "/var/html/html/assets/img/blank.jpg";
    const userName = "Hi, " + (localStorage.getItem("name") || "Guest");

    const dropdownMenu = document.querySelector(".dropdown-menu .dropdown-user-scroll");
    if (!dropdownMenu) {
        console.info("❌ Element .dropdown-menu .dropdown-user-scroll not found!");
        return;
    }

    const authItem = document.createElement("li");
    authItem.innerHTML = `<div class="dropdown-divider"></div>`;

    if (userEmail) {
        console.log("✅ User authenticated, showing Logout button.");

        // Update UI with user info
        const profileUsername = document.querySelector(".profile-username .fw-bold");
        const userBoxName = document.querySelector(".user-box .u-text h4");
        const userBoxEmail = document.querySelector(".user-box .u-text p");
        const topbarAvatar = document.querySelector(".profile-pic .avatar-img");
        const dropdownAvatar = document.querySelector(".user-box .avatar-lg img");

        if (profileUsername) profileUsername.textContent = userName;
        if (userBoxName) userBoxName.textContent = userName;
        if (userBoxEmail) userBoxEmail.textContent = userEmail;
        if (topbarAvatar) topbarAvatar.src = userPicture;
        if (dropdownAvatar) dropdownAvatar.src = userPicture;

        authItem.innerHTML += `<a class="dropdown-item" href="#" id="logout-btn">Logout</a>`;
    } else {
        console.log("🔴 User not authenticated, showing Login button.");

        const profilePic = document.querySelector(".dropdown-toggle.profile-pic");
        const profileUsername = document.querySelector(".profile-username .fw-bold");
        const avatarImg = document.querySelector(".avatar-img");

        dropdownMenu.classList.add("d-none");
        dropdownMenu.classList.add("disabled-dropdown");

        if (profilePic) {
            profilePic.innerHTML = `
                <button class="btn btn-primary" id="login-btn">
                    <i class="fab fa-google"></i> Login with Gmail
                </button>
            `;
        }

        if (profileUsername) profileUsername.textContent = "";
        if (avatarImg) avatarImg.src = "/var/html/html/assets/img/blank.jpg";

        const loginBtn = document.getElementById("login-btn");
        if (loginBtn) {
            loginBtn.addEventListener("click", () => {
                window.location.href = "login.html";
            });
        }
    }

    dropdownMenu.appendChild(authItem);

    setTimeout(() => {
        const logoutBtn = document.getElementById("logout-btn");
        if (logoutBtn) {
            logoutBtn.addEventListener("click", (event) => {
                event.preventDefault();
                swal({
                    title: "Logout",
                    text: "Are you sure you want to logout?",
                    icon: "warning",
                    buttons: {
                        confirm: {
                            text: "Yes, logout!",
                            className: "btn btn-success",
                        },
                        cancel: {
                            visible: true,
                            className: "btn btn-danger",
                        },
                    },
                }).then((willLogout) => {
                    if (willLogout) {
                        
                        localStorage.clear();
                         logoutGoogle(); // Logout dari Gmail juga
                        window.location.href = "forms_insert.html";
                    }
                });
            });
        } else {
            console.info("❌ Logout button not found!");
        }
    }, 100);
}


// ✅ Fungsi ambil token
export async function GetTokenGeneral() {
    try {
      console.log("🔍 Fetching Token...");
  
      const response = await fetch(CONFIG.URL_GET_TOKEN, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({
          grant_type: "client_credentials",
          client_id: CONFIG.CLIENT_ID,
          client_secret: CONFIG.CLIENT_SECRET,
          scope: "case"
        })
      });
  
      const data = await response.json();
      console.log("🔑 Get Token Response:", data);
  
      if (!response.ok || !data.access_token) {
        throw new Error(`Token Error:\n ${data.response_message || "Unknown error"}`);
      }
  
      return data.access_token;
    } catch (error) {
      console.error("❌ Error fetching token:", error);
      swal({ title: "Error", text: error.message, icon: "error" });
    }
  }
  
 
  function logoutGoogle() {
    const iframe = document.createElement('iframe');
    iframe.src = "https://accounts.google.com/Logout";
    iframe.style.display = "none";
    document.body.appendChild(iframe);

    setTimeout(() => {
        document.body.removeChild(iframe);
    }, 2000); // Tunggu logout selesai
}

export default CONFIG;
