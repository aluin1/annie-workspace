import CONFIG from './config_data_case.js';
import { GetTokenGeneral } from './config_data_case.js';

const URL_GET_DATA_ANNIE = CONFIG.URL_GET_DATA_ANNIE; 
const URL_INSERT_CASE = CONFIG.URL_INSERT_CASE; 
 
console.log("🔗 URL_GET_DATA_ANNIE:", URL_GET_DATA_ANNIE);
console.log("🔗 URL_INSERT_CASE:", URL_INSERT_CASE);

document.addEventListener('DOMContentLoaded', async function () {
  await fetchData(); // Panggil saat halaman dimuat
});

// ✅ Ambil TPID dari URL
function getTPID() {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get('tpid');
}

// ✅ Ambil isi input text
function getValue(id) {
  const el = document.getElementById(id);
  return el ? el.value : "";
}

// ✅ Ambil isi radio/checkbox
function getCheckedValue(name) {
  const el = document.querySelector(`input[name="${name}"]:checked`);
  return el ? el.value : "";
}

// ✅ Fungsi ambil data case dari Annie
async function fetchData() {
  const loadingIndicator = document.getElementById("loading");

  try {
    if (loadingIndicator) loadingIndicator.style.display = "block";

    const token = await GetTokenGeneral();

    const requestData = {
      TimePointId: getTPID()
    };

    console.log("📤 Requesting data with:", requestData);

    swal({ title: "Get Data Case From Annie...", text: "Please wait...", icon: "info", buttons: false });

    const dataResponse = await fetch(URL_GET_DATA_ANNIE, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(requestData)
    });

    const data = await dataResponse.json(); 
    
    if (!dataResponse.ok) {
      throw new Error(`Error Message: \n${data.response_message.replace(/;/g, '\n')}`);
  }


    if (!Array.isArray(data) || data.length === 0) {
      throw new Error('Data Case Empty From Annie');
    }

    const dataFirst = data[0];
    console.log("📥 Get Data First:", dataFirst);

    document.getElementById("customer_number").value = dataFirst["CustomerNumber"];
    document.getElementById("email").value = dataFirst["Email"];
    document.getElementById("comment").value = dataFirst["Comment"];

    swal.close();
  } catch (error) {
    console.error('❌ Error:', error);
    swal(error.message, { icon: "error", buttons: { confirm: { className: "btn btn-danger" } } });
  } finally {
    if (loadingIndicator) loadingIndicator.style.display = "none";
  }
}
 

  document.getElementById("submit_btn").addEventListener("click", async function (event) {
    event.preventDefault(); // ❗ Hindari reload halaman
  const token = await GetTokenGeneral();

  const requestDataInsert = {
    customer_number: getValue('customer_number'),  
    doctor_name: getValue('doctor_name'),
    email: getValue('email'),
    previous_case: getCheckedValue('previous_case'),
    previous_case_number: getValue('previous_case_number'),
    patient_name: getValue('patient_name'),
    dob: getValue('dob'),
    height_of_patient: getValue('height_of_patient'),
    gender: getCheckedValue('gender'),
    race: getValue('race'),
    package_list: getValue('package_list'),
    lateral_xray_date: getValue('lateral_xray_date'),
    consult_date: getValue('consult_date'),
    missing_teeth: getValue('missing_teeth'),
    adenoids_removed: getCheckedValue('adenoids_removed'),
    comment: getValue('comment'),
  };

  swal({ title: "Insert Data Case", text: "Please wait...", icon: "info", buttons: false });

  try {
    const response = await fetch(URL_INSERT_CASE, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(requestDataInsert),
    });

    const result = await response.json();
    if (!response.ok) {
        throw new Error(`Error Message: \n${result.response_message.replace(/;/g, '\n')}`);
    }


    console.log("✅ Response:", result);

    swal.close();
    swal({
      title: "Success",
      text: "Case submitted successfully!",
      icon: "success",
      button: "OK"
    }).then(() => {
      window.location.href = "data_case.html";
    });

  } catch (err) {
    console.error("❌ Error:", err);
    swal.close();
    swal({ title: "Error", text: err.message, icon: "error" }); 
  }
});
 