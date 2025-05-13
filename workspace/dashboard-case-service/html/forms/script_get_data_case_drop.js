import CONFIG from './config_data_case.js';
import { GetTokenGeneral } from './config_data_case.js';

const URL_GET_CASE = CONFIG.URL_GET_CASE; 
const URL_EDIT_STATUS_CASE = CONFIG.URL_EDIT_STATUS_CASE;
const URL_VALIDATION_TOKEN_GMAIL =  CONFIG.URL_VALIDATION_TOKEN_GMAIL;
 
window.addEventListener("pageshow", function (event) {
    if (event.persisted) {
        location.reload();
    }
});


document.addEventListener('DOMContentLoaded', async function() {  
    fetchData()
});

// ✅ **Fungsi Fetch Data dengan Logging**
async function fetchData() {
    try {
         
        swal({ title: "Loading Check Access...", text: "Please wait...", icon: "info", buttons: false });
        const tokenGmail = localStorage.getItem("google_jwt_token");
        // const tokenGmail =null;
        //console.log("📥 Google JWT Token:", tokenGmail);
        
    
            const validationToken = await ValidationToken(tokenGmail);
            console.log("validation Token Gmail: "+validationToken);

            if (validationToken){
                
             swal({ title: "Loading Get Data...", text: "Please wait...", icon: "info", buttons: false });
                const token = await GetTokenGeneral(); // 🔑 Ambil token sebelum submit
                console.log("✅ Token received successfully.");

                // console.log(`📡 Fetching case data from API: ${URL_GET_CASE}?page=${page}&page_size=${pageSize}`);
                // const response = await fetch(`${URL_GET_CASE}?page=${page}&page_size=${pageSize}`, {
                    console.log(`📡 Fetching case data from API: ${URL_GET_CASE}?page=&page_size=`);
                    const response = await fetch(`${URL_GET_CASE}?page=&page_size=`, {
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
            
        swal.close();
        });
            }
          


    } catch (error) {
        swal.close();
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
        const LateralXrayImage =  caseData.lateral_xray_image;
        const FrontalXrayImage = caseData.frontal_xray_image;
        const LowerArchImage= caseData.lower_arch_image;
        const UpperArchImage = caseData.upper_arch_image;
        const HandwristXrayImage = caseData.handwrist_xray_image;
        const PanoramicXrayImage = caseData.panoramic_xray_image;
        const AdditionalRecord1 = caseData.additional_record_1;
        const AdditionalRecord2 = caseData.additional_record_2;
        const AdditionalRecord3 = caseData.additional_record_3;
        const AdditionalRecord4 = caseData.additional_record_4;
        const AdditionalRecord5 = caseData.additional_record_5;

        let detailHtml = `
            <h6><b>Completed: <br>${dateCompleted}</b></h6>
            <h6><b>Annie VIP Order Form: </b></h6>
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
                    <td style="width: 200;" ><strong>Lateral X-Ray Image:</strong></td><td><img src="${LateralXrayImage}" style="${LateralXrayImage ? '' : 'display:none'}"  width="300">
                    <br><a href ="${LateralXrayImage}" target="_blank"  style="${LateralXrayImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a>
                    </td> 


                    <td style="width: 200;"><strong>Frontal X-Ray Image:</strong></td><td><img src="${FrontalXrayImage}" style="${FrontalXrayImage ? '' : 'display:none'}" width="300">
                    <br><a href ="${FrontalXrayImage}" target="_blank" style="${FrontalXrayImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 


                    </tr> 
                    <tr>


                    <td><strong>Lower Arch Image:</strong></td><td><img src="${LowerArchImage}" width="300" style="${LowerArchImage ? '' : 'display:none'}" >
                    <br><a href ="${LowerArchImage}" target="_blank" style="${LowerArchImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 

                    
					<td><strong>Upper Arch Image:</strong></td><td><img src="${UpperArchImage}" width="300" style="${UpperArchImage ? '' : 'display:none'}" >
                    <br><a href ="${UpperArchImage}" target="_blank" style="${UpperArchImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 


                    </tr> 
                    <tr>


                    <td><strong>HandWrist X-Ray Image:</strong></td><td><img src="${HandwristXrayImage}" width="300" style="${HandwristXrayImage ? '' : 'display:none'}" >
                    <br><a href ="${HandwristXrayImage}" target="_blank" style="${HandwristXrayImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 

 
                    <td><strong>Panoramic Xray Image:</strong></td><td><img src="${PanoramicXrayImage}" width="300" style="${PanoramicXrayImage ? '' : 'display:none'}" >
                    <br><a href ="${PanoramicXrayImage}" target="_blank" style="${PanoramicXrayImage ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 


                    </tr> 
                    <tr>

                    <td><strong>Additional Record 1:</strong></td><td><img src="${AdditionalRecord1}" width="300" style="${AdditionalRecord1 ? '' : 'display:none'}" >
                    <br><a href ="${AdditionalRecord1}" target="_blank" style="${AdditionalRecord1 ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 

                    <td><strong>Additional Record 2:</strong></td><td><img src="${AdditionalRecord2}" width="300" style="${AdditionalRecord2 ? '' : 'display:none'}" >
                    <br><a href ="${AdditionalRecord2}" target="_blank" style="${AdditionalRecord2 ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 
 


                    </tr> 
                    <tr>


                    <td><strong>Additional Record 3:</strong></td><td><img src="${AdditionalRecord3}" width="300" style="${AdditionalRecord3 ? '' : 'display:none'}" >
                    <br><a href ="${AdditionalRecord3}" target="_blank" style="${AdditionalRecord3 ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 

                    <td><strong>Additional Record 4:</strong></td><td><img src="${AdditionalRecord4}" width="300" style="${AdditionalRecord4 ? '' : 'display:none'}" >
                    <br><a href ="${AdditionalRecord4}" target="_blank" style="${AdditionalRecord4 ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 
 


                    </tr> 
                    <tr>


                    <td><strong>Additional Record 5:</strong></td><td><img src="${AdditionalRecord5}" width="300" style="${AdditionalRecord5 ? '' : 'display:none'}" >
                    <br><a href ="${AdditionalRecord5}" target="_blank" style="${AdditionalRecord5 ? '' : 'display:none'}" ><button class="btn btn-info mt-3">View</button></a> 
                    </td> 
 

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

            const token = await GetTokenGeneral(); // 🔑 Ambil token sebelum submit
            const requestBody = {
                customer_number: caseData.customer_number,
                case_id: String(caseData.case_id),
                status_case: newStatus,
                comment: newComment
            };

            console.log("🚀 Sending update request", requestBody);

            try {
                
        swal({ title: "Loading Update Data Case...", text: "Please wait...", icon: "info", buttons: false });
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
                
                swal.close();
                swal({
                    title: "Success",
                    text: "Update Data Case successfully!",
                    icon: "success",
                    button: "OK" // Tombol konfirmasi
                }).then(() => { 
                    
                modal.hide();
                fetchData(); // Refresh data after update
                });
 
            } catch (error) {
                console.error("❌ Error updating status:", error);
                swal.close();
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
 
const statusColors = {
    "1": "bg-primary",    
    "2": "bg-secondary",  
    "3": "bg-success"   ,  
    "4": "bg-info"    ,  
    "5": "bg-warning"     
};

async function ValidationToken(tokenGmail) {
   
    if (tokenGmail == null ||tokenGmail =="" ){
             
        swal.close();
        swal({
            title: "Please Login First",
            text: "You need to log in before continue Process.",
            icon: "error" ,
            button: "OK" 
        }).then(() => {
            window.location.href = "login.html"; // Redirect setelah klik OK
        });
        return false;
        
    }  else {

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

            // await swal({
            //     title: "Success",
            //     text: `Welcome, ${data.name || "User"}!`,
            //     icon: "success",
            //     button: "OK"
            // });

            // window.location.href = "data_case.html"; // Redirect setelah klik OK
            return true; // Token valid
        } else {
            console.error("resp validate token auth:", data.response_message);
            throw new Error(data.response_message);
        }
    } catch (err) {
        console.error("❌ Error validate token auth:", err);
        // swal.close();

        await swal({
            title: "Credential Failed",
            text: err.message,
            icon: "error",
            buttons: { confirm: { className: "btn btn-danger" } }
        });

        fetchData();
        window.location.href = "login.html";
        return false; // Token tidak valid
    }
}
}
