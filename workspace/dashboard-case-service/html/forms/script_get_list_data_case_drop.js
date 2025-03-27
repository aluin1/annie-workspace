import CONFIG from './config_data.js';

const URL_GET_CASE = CONFIG.URL_GET_CASE;
const URL_GET_TOKEN = CONFIG.URL_GET_TOKEN;
const URL_EDIT_STATUS_CASE = CONFIG.URL_EDIT_STATUS_CASE;
const CLIENT_ID = CONFIG.CLIENT_ID;
const CLIENT_SECRET = CONFIG.CLIENT_SECRET;

document.addEventListener('DOMContentLoaded', async function() {
    console.log("🚀 DOM Loaded - Fetching data...");
    fetchData();
});

// ✅ **Fungsi Fetch Data dengan Logging**
async function fetchData( pageSize, page) {
    try {
         
        const token = await getToken(); // 🔑 Ambil token sebelum submit
        console.log("✅ Token received successfully.");

        console.log(`📡 Fetching case data from API: ${URL_GET_CASE}?page=${page}&page_size=${pageSize}`);
        const response = await fetch(`${URL_GET_CASE}?page=${page}&page_size=${pageSize}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        const result = await response.json();
        if (!response.ok) throw new Error(result.response_message);
        console.log(`📄 Total cases received: ${result.data_case.length}`);

        // Kosongkan elemen sebelum menambahkan data baru
        console.log("🧹 Clearing old case data...");
        ["Submited-Case", "WorkInProgress-Case", "Completed-Case", "Abandoned-Case", "Trash-Case"].forEach(id => {
            document.getElementById(id).innerHTML = "";
        });

        // Pisahkan berdasarkan status_case
        result.data_case.forEach((caseItem, index) => {
            const textName = caseItem.doctor_name + " | " + caseItem.email;
            const initials = getInitials(caseItem.doctor_name);
            // console.log(`🔍 Processing case #${index + 1}: ${caseItem.case_id} (Doctor: ${caseItem.doctor_name})`);
            let caseStatusClass = statusColors[caseItem.status_case] || "bg-secondary";

            let caseHtml = `
            <div class="card p-3"  data-id='${JSON.stringify(caseItem)}'>
                <div class="d-flex align-items-center ">
                    <span class="stamp ${caseStatusClass} me-2" style="color: #fff;">${initials}</span>
                    <div>
                        <b> <small>#${caseItem.case_id}</small> </b>
                        <div class="float-end text-primary"><small>${timeAgo(caseItem.time_create)}</small></div><br>
                        <small class="text-muted">${truncateText(textName, 25)}</small>
                    </div>
                </div> </div>
            `;

            if (caseItem.status_case === "1") {
                document.getElementById("Submited-Case").innerHTML += caseHtml;
            } else if (caseItem.status_case === "2") {
                document.getElementById("WorkInProgress-Case").innerHTML += caseHtml;
            } else if (caseItem.status_case === "3") {
                document.getElementById("Completed-Case").innerHTML += caseHtml;
            } else if (caseItem.status_case === "4") {
                document.getElementById("Abandoned-Case").innerHTML += caseHtml;
            } else if (caseItem.status_case === "5") {
                document.getElementById("Trash-Case").innerHTML += caseHtml;
            }

            // console.log(`✅ Case #${caseItem.case_id} added to category ${caseItem.status_case}`);
        });

    } catch (error) {
        console.error('❌ Error:', error.message);
        swal(error.message, { icon: "error", buttons: { confirm: { className: "btn btn-danger" } } });
    }
}

// ✅ **Fungsi Format Waktu "a minute ago"**
function timeAgo(timestamp) {
    const now = new Date();
    const past = new Date(timestamp);
    const diffInSeconds = Math.floor((now - past) / 1000);

    if (diffInSeconds < 60) return "a minute ago";
    else if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    else if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    else return `${Math.floor(diffInSeconds / 86400)} days ago`;
}

// ✅ **Fungsi Ambil Inisial Nama**
function getInitials(name) {
    if (!name) return "";  // Jika nama kosong, kembalikan string kosong
    return name
        .split(" ")            // Pecah berdasarkan spasi
        .map(word => word[0])  // Ambil huruf pertama dari setiap kata
        .join("")              // Gabungkan hasilnya
        .toUpperCase();        // Pastikan huruf kapital
}

// ✅ **Fungsi Memotong Teks**
function truncateText(text, maxLength) {
    if (!text) return "";
    return text.length > maxLength ? text.slice(0, maxLength) + "..." : text;
}
document.addEventListener("click", function(event) {
    const card = event.target.closest(".card.p-3");
    if (!card) return; // Jika tidak klik di dalam kartu, keluar

    console.log("📋 Card clicked:", card);

    try {
        const caseData = JSON.parse(card.getAttribute("data-id"));
        console.log("📄 Case Data:", caseData);

        const dateCompleted = formatDate(caseData.time_create);
        let detailHtml = `
            <h6><b>Completed: <br>${dateCompleted}</b></h6>
            <h6><b>Biogresive Order Form - RMODS: </b></h6>
            <table border='0' class='display table table-striped table-hover'>
                <tr>
                    <td><strong>Customer Number:</strong></td><td> ${caseData.customer_number}</td>                    
                    <td><strong>Gender:</strong></td><td> ${caseData.gender}</td>
                </tr> 
                <tr>
                    <td><strong>Doctor Name:</strong></td><td> ${caseData.doctor_name}</td>
                    <td><strong>Race:</strong></td><td> ${caseData.race}</td>
                </tr> 
                <tr>
                    <td><strong>Email:</strong></td><td> ${caseData.email}</td>
                    <td><strong>Package List:</strong></td><td> ${caseData.package_list}</td>
                </tr> 
                <tr>
                    <td><strong>Previous Case Number:</strong></td><td> ${caseData.previous_case_number}</td>
                    <td><strong>Consult Date:</strong></td><td>${formatDateOnly(caseData.consult_date)}</td>
                </tr> 
                <tr>
                    <td><strong>Patient Name:</strong></td><td> ${caseData.patient_name}</td>
                    <td><strong>Missing Teeth:</strong></td><td> ${caseData.missing_teeth}</td>
                </tr> 
                <tr>
                    <td><strong>DOB:</strong></td><td> ${formatDateOnly(caseData.dob)}</td>
                    <td><strong>Adenoids Removed:</strong></td><td> ${caseData.adenoids_removed}</td>
                </tr> 
                    <tr>
                    <td><strong>Height:</strong></td><td> ${caseData.height_of_patient}</td>
                     <td><strong>Lateral X-ray Date</strong></td><td> ${formatDateOnly(caseData.lateral_xray_date)}</td>
                    </tr> 
               
                    <tr>
                    <td><strong>Lateral X-Ray Image:</strong></td><td> ${caseData.lateral_xray_image}</td> 
                    <td><strong>Frontal X-Ray Image:</strong></td><td> ${caseData.frontal_xray_image}</td>
                    </tr> 
                    <tr>
                    <td><strong>Lower Arch Image:</strong></td><td> ${caseData.lower_arch_image}</td>
					<td><strong>Upper Arch Image:</strong></td><td> ${caseData.upper_arch_image}</td>
                    </tr> 
                    <tr>
                    <td><strong>HandWrist X-Ray Image:</strong></td><td> ${caseData.handwrist_xray_image}</td>
					<td><strong>Panoramic X-Ray (Panorex) Image:</strong></td><td> ${caseData.panoramic_xray_image}</td>
                    </tr> 
                    <tr>
                    <td><strong>Additional Record 1:</strong></td><td> ${caseData.additional_record_1}</td>
					<td><strong>Additional Record 2:</strong></td><td> ${caseData.additional_record_2}</td>
                    </tr> 
                    <tr>
                    <td><strong>Additional Record 3:</strong></td><td> ${caseData.additional_record_3}</td>
					<td><strong>Additional Record 4:</strong></td><td> ${caseData.additional_record_4}</td>
                    </tr> 
                    <tr>
                    <td><strong>Additional Record 5:</strong></td><td> ${caseData.additional_record_5}</td>
					   <td></td><td></td>
                    </tr> 
                <tr>
                    <td colspan='1'><strong>Comment:</strong></td>
                    <td colspan='3'>
                    <div class="form-group"> 
                        <textarea id="commentInput" class="form-control" rows="3">${caseData.comment}</textarea>
                        </div>

                    </td>
                </tr>  
                 <tr>
                    <td><strong>Move Status:</strong></td>
                    <td colspan='3'>
                        <select id="statusSelect" class="form-select form-control">
                            <option value="">Move to</option> 
                             <option value="1" ${caseData.status_case === "1" ? "selected" : ""}>Submitted</option> 
                                <option value="2" ${caseData.status_case === "2" ? "selected" : ""}>Work In Progress</option> 
                                <option value="3" ${caseData.status_case === "3" ? "selected" : ""}>Completed</option> 
                                <option value="4" ${caseData.status_case === "4" ? "selected" : ""}>Abandoned</option> 
                                <option value="5" ${caseData.status_case === "5" ? "selected" : ""}>Trash</option> 
                            </select>
                    </td>
                </tr>
            </table>
            <div class="modal-footer"> 
            <button type="button"  id="saveStatusBtn" class="btn btn-info mt-3">Save</button>
                  <button type="button" class="btn btn-danger  mt-3" data-bs-dismiss="modal">Close</button>
              </div>
        `;

        document.querySelector("#modalBody").innerHTML = detailHtml;
        document.getElementById("customerNumber").textContent = "#"+caseData.case_id;

        console.log("🛠️ Showing modal...");
        const modal = new bootstrap.Modal(document.getElementById("detailModal"));
        modal.show();

        // Handle Save Button Click
        document.getElementById("saveStatusBtn").addEventListener("click", async function () {
            const newStatus = document.getElementById("statusSelect").value;
            const newComment = document.getElementById("commentInput").value;

            if (!newStatus) { 
                swal({ title: "Error", text: "Please select a status", icon: "error" });
                return;
            }

            const token = await getToken(); // 🔑 Ambil token sebelum submit
            const requestBody = {
                customer_number: caseData.customer_number,
                case_id: String(caseData.case_id),
                status_case: newStatus,
                comment: newComment
            };

            console.log("🚀 Sending update request", requestBody);

            try {
                const response = await fetch(URL_EDIT_STATUS_CASE , {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(requestBody)
                });
                
                const result = await response.json();
                if (!response.ok) throw new Error(result.response_message);
                
                console.log("✅ Status updated successfully", result); 
                
                swal({
                    title: "Success",
                    text: "Update Data Case successfully!",
                    icon: "success",
                    button: "OK" // Tombol konfirmasi
                }).then(() => {
                    window.location.href = "data_case.html"; // Redirect setelah klik OK
                });

                modal.hide();
                fetchData(); // Refresh data after update
            } catch (error) {
                console.error("❌ Error updating status:", error);
                swal({ title: "Error", text: error.message, icon: "error" });
            }
        });

    } catch (error) {
        console.error("❌ Error parsing case data:", error);
        swal({ title: "Error", text: error.message, icon: "error" });
    }
});

// ✅ **Fungsi Format Tanggal**
function formatDateOnly(dateString) {
    const date = new Date(dateString.replace(" ", "T")); // Pastikan format ISO
    return date.toLocaleString("en-GB", {
        day: "2-digit",
        month: "long",
        year: "numeric",
    });
}
// ✅ **Fungsi Format Tanggal**
function formatDate(dateString) {
    const date = new Date(dateString.replace(" ", "T")); // Pastikan format ISO
    return date.toLocaleString("en-GB", {
        day: "2-digit",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
        hour12: true,
    });
}

async function getToken() {
    try {
        console.log("🔍 Fetching Token...");

        const response = await fetch(URL_GET_TOKEN, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: new URLSearchParams({
                grant_type: "client_credentials",
                client_id: CLIENT_ID,
                client_secret: CLIENT_SECRET,
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

const statusColors = {
    "1": "bg-primary",    
    "2": "bg-secondary",  
    "3": "bg-success"   ,  
    "4": "bg-info"    ,  
    "5": "bg-warning"     
};
