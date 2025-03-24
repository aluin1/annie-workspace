import CONFIG from './config_data.js';

const URL_GET_CASE = CONFIG.URL_GET_CASE;
const URL_GET_TOKEN = CONFIG.URL_GET_TOKEN;
const CLIENT_ID = CONFIG.CLIENT_ID;
const CLIENT_SECRET = CONFIG.CLIENT_SECRET;

console.log("🔗 URL_GET_TOKEN:", URL_GET_TOKEN);
console.log("🔗 URL_GET_CASE:", URL_GET_CASE);

document.addEventListener('DOMContentLoaded', async function() {
    fetchData(50,1); // Ambil data pertama kali dengan page_size = 50
});

// ✅ **Fungsi Fetch Data**
async function fetchData(pageSize ,page ) {
    const tbody = document.querySelector("#multi-filter-select tbody");
    const loadingIndicator = document.getElementById("loading");

    try {
        loadingIndicator.style.display = "block"; // Tampilkan loading
        console.log(`🔄 Fetching data with pageSize: ${pageSize}, page: ${page}`);

        // 1️⃣ **Get Token**
        const formData = new URLSearchParams();
        formData.append("grant_type", "client_credentials");
        formData.append("client_id", CLIENT_ID);
        formData.append("client_secret", CLIENT_SECRET);
        formData.append("scope", "case");

        console.log("📡 Requesting Token...");
        const tokenResponse = await fetch(URL_GET_TOKEN, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: formData
        });

        const tokenData = await tokenResponse.json();
        console.log("🔑 Token Response:", tokenData);

        if (!tokenResponse.ok) throw new Error(`Get token failed: ${tokenData.response_message}`);

        const token = tokenData.access_token;
        if (!token) throw new Error("Access token is missing.");

        console.log("✅ Access token obtained successfully");

        // 2️⃣ **Get Data with Token**
        const urlWithParams = `${URL_GET_CASE}?page=${page}&page_size=${pageSize}`;
        console.log("📡 Requesting Data from:", urlWithParams);

        const response = await fetch(urlWithParams, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        const result = await response.json();
        console.log("📥 API Response:", result);

        if (!response.ok) throw new Error(`API Error: ${result.response_message}`);

        console.log("✅ Data fetched successfully:", result.data_case.length, "records");

        // 3️⃣ **Populate Table**
        tbody.innerHTML = ""; // Kosongkan tabel sebelum mengisi data
        
        // 🔽 **Sort data berdasarkan case_id terbesar** 🔽
        result.data_case.sort((a, b) => b.case_id - a.case_id);         
        console.log("📊 Sorted case IDs:", result.data_case.map(c => c.case_id));
        console.log("📥 API Response Data Case:", result.data_case);

        result.data_case.forEach((caseItem) => {
            let row = document.createElement("tr");
            row.innerHTML = `
                <td>${caseItem.case_id} - ${caseItem.customer_number}</td>
                <td>${caseItem.doctor_name}</td>
                <td>${caseItem.email}</td>
                <td>${caseItem.patient_name}</td>
                <td>${caseItem.dob}</td>
                <td>${caseItem.gender}</td>
                <td>
                    <button class="btn btn-info btn-sm detail-btn" data-id='${JSON.stringify(caseItem)}'>
                        Detail
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });

        console.log(`✅ Table updated with ${result.data_case.length} records`);

        // 4️⃣ **Hapus DataTables Lama & Inisialisasi Ulang**
        if ($.fn.DataTable.isDataTable("#multi-filter-select")) {
            $("#multi-filter-select").DataTable().destroy();
        }

        $("#multi-filter-select").DataTable({
            pageLength: 10, // Gunakan pageSize dari parameter
            lengthMenu: [5, 10, 25, 50], // Dropdown Show Entries
            order: [], // Ini memastikan tidak ada sorting default dari DataTable
            initComplete: function () {
                this.api().columns().every(function () {
                    var column = this;
                    var select = $('<select class="form-select"><option value=""></option></select>')
                        .appendTo($(column.footer()).empty())
                        .on("change", function () {
                            var val = $.fn.dataTable.util.escapeRegex($(this).val());
                            column.search(val ? "^" + val + "$" : "", true, false).draw();
                        });

                    column.data().unique().sort().each(function (d, j) {
                        select.append('<option value="' + d + '">' + d + "</option>");
                    });
                });
            }
        });

        console.log("✅ DataTable initialized with pageSize:", pageSize);

    } catch (error) {
        console.error('❌ Error:', error.message);
        swal(error.message, { icon: "error", buttons: { confirm: { className: "btn btn-danger" } } });
    } finally {
        loadingIndicator.style.display = "none"; // Pastikan loading di-hide walaupun error
    }
}

// ✅ **Event Listener untuk Perubahan "Show Entries"**
$('#multi-filter-select').on('length.dt', function (e, settings, len) {
    console.log(`🔄 Show Entries changed to: ${len}. Fetching new data...`);
    fetchData(len); // Hit API dengan page_size baru
});

// ✅ **Event Delegation untuk Tombol Detail**
document.addEventListener("click", function(event) {
    if (event.target.classList.contains("detail-btn")) {
        const caseData = JSON.parse(event.target.getAttribute("data-id"));
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
                    <td><strong>Consult Date:</strong></td><td> ${caseData.consult_date}</td>
                </tr> 
                <tr>
                    <td><strong>Patient Name:</strong></td><td> ${caseData.patient_name}</td>
                    <td><strong>Missing Teeth:</strong></td><td> ${caseData.missing_teeth}</td>
                </tr> 
                <tr>
                    <td><strong>DOB:</strong></td><td> ${caseData.dob}</td>
                    <td><strong>Adenoids Removed:</strong></td><td> ${caseData.adenoids_removed}</td>
                </tr> 
                    <tr>
                    <td><strong>Height:</strong></td><td> ${caseData.height_of_patient}</td>
                     <td></td><td></td>
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
                    <td colspan='3'>${caseData.comment}</td>
                </tr>  
            </table>
        `;

        document.querySelector("#modalBody").innerHTML = detailHtml;
        document.getElementById("customerNumber").textContent = caseData.customer_number;
        $("#detailModal").modal("show");

        console.log("📋 Case details opened for customer:", caseData.customer_number);
    }
});

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
