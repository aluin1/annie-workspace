// const CONFIG = {
//     URL_INSERT_CASE: 'http://localhost:8081/insertCase',
//     URL_GET_CASE: 'http://localhost:8081/dataCase',
//     URL_GET_TOKEN: 'http://localhost:8081/getToken', 
//     URL_EDIT_STATUS_CASE: 'http://localhost:8081/editCaseStatus', 
//     CLIENT_ID:'79d4270a98d241999fa7e54e015b51ab',
//     CLIENT_SECRET:'22106803111c4cdda620875ed4779129' ,
//     URL_UPLOAD_FILE:'http://localhost:8081/uploadFiles',
//     EMAIL_SUPPORT:'fmbods@bioprogressive.org',
//     URL_VALIDATION_TOKEN_GMAIL:'http://localhost:8081/validationTokenGmail'
// };

const CONFIG = {
    URL_INSERT_CASE: 'https://annievipbe.ismartdds.com/insertCase',
    URL_GET_CASE: 'https://annievipbe.ismartdds.com/dataCase',
    URL_GET_TOKEN: 'https://annievipbe.ismartdds.com/getToken', 
    URL_EDIT_STATUS_CASE: 'https://annievipbe.ismartdds.com/editCaseStatus', 
    CLIENT_ID:'79d4270a98d241999fa7e54e015b51ab',
    CLIENT_SECRET:'22106803111c4cdda620875ed4779129' ,
    URL_UPLOAD_FILE:'https://annievipbe.ismartdds.com/uploadFiles',
    EMAIL_SUPPORT:'fmbods@bioprogressive.org',
    URL_VALIDATION_TOKEN_GMAIL:'http://localhost:8081/validationTokenGmail'
};

document.addEventListener("DOMContentLoaded", checkDataUser);

async function checkDataUser() {
    console.log("🚀 DOM Loaded - Checking authentication...");

    let userEmail = localStorage.getItem("email") || "";
    let userPicture = localStorage.getItem("picture") || "/assets/img/blank.jpg";
    let userName = "Hi, " + (localStorage.getItem("name") || "Guest");

    let dropdownMenu = document.querySelector(".dropdown-menu .dropdown-user-scroll");
    if (!dropdownMenu) {
        console.error("❌ Element .dropdown-menu .dropdown-user-scroll not found!");
        return;
    }

    let authItem = document.createElement("li");
    authItem.innerHTML = `<div class="dropdown-divider"></div>`;

    if (userEmail) {
        console.log("✅ User authenticated, showing Logout button.");

        // Set profile-related DOM elements
        let profileUsername = document.querySelector(".profile-username .fw-bold");
        let userBoxName = document.querySelector(".user-box .u-text h4");
        let userBoxEmail = document.querySelector(".user-box .u-text p");
        let topbarAvatar = document.querySelector(".profile-pic .avatar-img");
        let dropdownAvatar = document.querySelector(".user-box .avatar-lg img");

        if (profileUsername) profileUsername.textContent = userName;
        if (userBoxName) userBoxName.textContent = userName;
        if (userBoxEmail) userBoxEmail.textContent = userEmail;
        if (topbarAvatar) topbarAvatar.src = userPicture;
        if (dropdownAvatar) dropdownAvatar.src = userPicture;

        authItem.innerHTML += `<a class="dropdown-item" href="#" id="logout-btn">Logout</a>`;
    } else {
        // Ganti menjadi tombol login dengan icon Gmail
        console.log("🔴 User not authenticated, showing Login button.");

        let profilePic = document.querySelector(".dropdown-toggle.profile-pic");
        let profileUsername = document.querySelector(".profile-username .fw-bold");
        let avatarImg = document.querySelector(".avatar-img");
        // dropdownMenu.classList.remove("d-none"); // Ensure dropdown is visible
        // dropdownMenu.classList.remove("disabled-dropdown"); // Enable dropdown menu
        
    // let dropdownMenuUser = document.querySelector(".user-dropdown .dropdown-menu .dropdown-user-scroll");
    
    //         dropdownMenuUser.classList.add("d-none"); // Hide the dropdown menu
    //         dropdownMenuUser.classList.add("disabled-dropdown"); // Disable dropdown interaction

        dropdownMenu.classList.add("d-none"); // Hide the dropdown menu
        dropdownMenu.classList.add("disabled-dropdown"); // Disable dropdown interaction

        if (profilePic) {
            profilePic.innerHTML = `
            <style>
            .navbar-expand-lg .navbar-nav .dropdown-menu {
                    display: none;
                    left: auto;
                    right: 0;
                    z-index: 1001;
                }
            </style>
                <button class="btn btn-primary" id="login-btn">
                    <i class="fab fa-google"></i> Login with Gmail
                </button>
            `;
        }
        
        if (profileUsername) profileUsername.textContent = ""; // Clear any text for logged-in user
        if (avatarImg) avatarImg.src = "/assets/img/blank.jpg"; // Placeholder image

        // Add event listener for the login button
        let loginBtn = document.getElementById("login-btn");
        if (loginBtn) {
            loginBtn.addEventListener("click", function () {
                window.location.href = "login.html"; // Redirect to login.html immediately
            });
        }
    }

    dropdownMenu.appendChild(authItem);

    // 🚀 Tunggu hingga elemen logout tersedia di DOM
    setTimeout(() => {
        let logoutBtn = document.getElementById("logout-btn");
        if (logoutBtn) {
            logoutBtn.addEventListener("click", function (event) {
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
                        window.location.href = "login.html";
                    }
                });
            });
        } else {
            console.error("❌ Logout button not found!");
        }
    }, 100); // Tunggu 100ms untuk memastikan elemen ada di DOM
}


export default CONFIG;
