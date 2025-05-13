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

    document.getElementById("customer_number_TEXT").textContent = formatTo7Digits(dataFirst["CustomerNumber"]);
    document.getElementById("customer_number").value = formatTo7Digits(dataFirst["CustomerNumber"]);

    
    document.getElementById("doctor_name_TEXT").textContent = dataFirst["DoctorName"];
    document.getElementById("doctor_name").value =  dataFirst["DoctorName"];

    document.getElementById("email_TEXT").textContent =  dataFirst["Email"];
    document.getElementById("email").value = dataFirst["Email"];

   
    document.getElementById("previous_case_number_TEXT").textContent = dataFirst["CustomerNumber"];
    document.getElementById("previous_case_number").value = dataFirst["CustomerNumber"];

    
    document.getElementById("patient_name_TEXT").textContent = dataFirst["PatientName"];
    document.getElementById("patient_name").value = dataFirst["PatientName"];
    
    document.getElementById("dob_TEXT").textContent =formatDateOnly( toInputDateValue(dataFirst["PatientBirthDate"]));
    document.getElementById("dob").value = toInputDateValue(dataFirst["PatientBirthDate"]);;
    
    document.getElementById("height_of_patient_TEXT").textContent = dataFirst["PatientHeight"];
    document.getElementById("height_of_patient").value = dataFirst["PatientHeight"]; 

    

    document.getElementById("race_TEXT").textContent = dataFirst["PatientRace"];
    document.getElementById("race").value =dataFirst["PatientRace"]; 

    document.getElementById("package_list_TEXT").textContent = "D1 Comprehensive";
    document.getElementById("package_list").value = "D1 Comprehensive";
    
    
    document.getElementById("lateral_xray_date_TEXT").textContent =formatDateOnly( toInputDateValue(dataFirst["Image_Lateral_Date"]));
    document.getElementById("lateral_xray_date").value =toInputDateValue(dataFirst["Image_Lateral_Date"]);

    document.getElementById("consult_date_TEXT").textContent = formatDateOnly(toInputDateValue(dataFirst["TimePointDate"]));
    document.getElementById("consult_date").value =  toInputDateValue(dataFirst["TimePointDate"]);
    
    document.getElementById("missing_teeth_TEXT").textContent =  dataFirst["IsMissingTeeth"];
    document.getElementById("missing_teeth").value = dataFirst["IsMissingTeeth"];

    
   
    // var radioValue = "Yes";
    // selectRadioByValue("previous_case", radioValue); 
    document.getElementById("previous_case_TEXT").textContent =  "Yes";
    document.getElementById("previous_case").value = "Yes";


    var radioAdenoidsRemovedValue = TranslateYesNo(dataFirst["PatientAdenoid"]); 
    // selectRadioByValue("adenoids_removed", radioAdenoidsRemovedValue);     
    document.getElementById("adenoids_removed_TEXT").textContent =  radioAdenoidsRemovedValue;
    document.getElementById("adenoids_removed").value = radioAdenoidsRemovedValue;

    
    var radioGenderValue = TranslateGender(dataFirst["PatientGender"]);  
    // selectRadioByValue("gender", radioGenderValue);  
    document.getElementById("gender_TEXT").textContent =  radioGenderValue;
    document.getElementById("gender").value = radioGenderValue;


    
    document.getElementById("comment_TEXT").textContent =  dataFirst["Comment"];
    document.getElementById("comment").value = dataFirst["Comment"];

    
    const noImage="https://placehold.co/600x400?text=No+Image";
    
   const Image_Lateral= dataFirst["Image_Lateral"];
    document.getElementById("lateral_xray_image_IMG").src = Image_Lateral  ? Image_Lateral  : noImage;
    document.getElementById("lateral_xray_image").value = Image_Lateral;

    const Image_Frontal= dataFirst["Image_Frontal"];
    document.getElementById("frontal_xray_image_IMG").src = Image_Frontal  ? Image_Frontal  : noImage;
    document.getElementById("frontal_xray_image").value =Image_Frontal;


    const Image_Lower= dataFirst["Image_Lower"];
    document.getElementById("lower_arch_image_IMG").src =  Image_Lower  ? Image_Lower  : noImage;
    document.getElementById("lower_arch_image").value = Image_Lower;


    const Image_Upper= dataFirst["Image_Upper"];
    document.getElementById("upper_arch_image_IMG").src =  Image_Upper  ? Image_Upper  : noImage;
    document.getElementById("upper_arch_image").value =Image_Upper;

    //belum ada mapping
    document.getElementById("handwrist_xray_image_IMG").src = noImage;
    document.getElementById("handwrist_xray_image").value = noImage;


   const Image_Panoramic= dataFirst["Image_Panoramic"];

    document.getElementById("panoramic_xray_image_IMG").src = Image_Panoramic  ? Image_Panoramic  : noImage;
    document.getElementById("panoramic_xray_image").value = Image_Panoramic;

    
    const Image_Profile= dataFirst["Image_Profile"];
    document.getElementById("additional_record_1_IMG").src =  Image_Profile  ? Image_Profile  : noImage;
    document.getElementById("additional_record_1").value = Image_Profile
    
    const Image_ProfileSmile= dataFirst["Image_ProfileSmile"];
    document.getElementById("additional_record_2_IMG").src = Image_ProfileSmile  ? Image_ProfileSmile  : noImage;
    document.getElementById("additional_record_2").value =  Image_ProfileSmile;

    const Image_ProfileNoSmile= dataFirst["Image_ProfileNoSmile"];
    document.getElementById("additional_record_3_IMG").src =  Image_ProfileNoSmile  ? Image_ProfileNoSmile  : noImage;
    document.getElementById("additional_record_3").value =  Image_ProfileNoSmile;
    
    //belum ada mapping
    document.getElementById("additional_record_4_IMG").src = noImage;
    document.getElementById("additional_record_4").value = noImage;

    //belum ada mapping
    document.getElementById("additional_record_5_IMG").src = noImage;
    document.getElementById("additional_record_5").value = noImage;

    

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
    // previous_case: getCheckedValue('previous_case'),
    previous_case: getValue('previous_case'),
    previous_case_number: getValue('previous_case_number'),
    patient_name: getValue('patient_name'),
    dob: getValue('dob'),
    height_of_patient: getValue('height_of_patient'),
    // gender: getCheckedValue('gender'),
    gender: getValue('gender'),
    race: getValue('race'),
    package_list: getValue('package_list'),
    lateral_xray_date: getValue('lateral_xray_date'),
    consult_date: getValue('consult_date'),
    missing_teeth: getValue('missing_teeth'),
    // adenoids_removed: getCheckedValue('adenoids_removed'),
    adenoids_removed: getValue('adenoids_removed'),
    comment: getValue('comment'),
 
            lateral_xray_image: getValue('lateral_xray_image'),
            frontal_xray_image: getValue('frontal_xray_image'),
            lower_arch_image: getValue('lower_arch_image'),
            upper_arch_image: getValue('upper_arch_image'),
            handwrist_xray_image: getValue('handwrist_xray_image'),
            panoramic_xray_image: getValue('panoramic_xray_image'),
            additional_record_1: getValue('additional_record_1'),
            additional_record_2: getValue('additional_record_2'),
            additional_record_3: getValue('additional_record_3'),
            additional_record_4: getValue('additional_record_4'),
            additional_record_5: getValue('additional_record_5'),
    
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

function formatTo7Digits(value) {
  // Ubah ke string dulu (jaga-jaga kalau input berupa number)
  let str = value.toString();

  if (str.length >= 7) {
    // Ambil 7 digit terakhir
    return str.slice(-7);
  } else {
    // Tambah 0 di depan hingga 7 digit
    return str.padStart(7, '0');
  }
}

function toInputDateValue(isoDateString) {
  const date = new Date(isoDateString);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // Bulan dimulai dari 0
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function selectRadioByValue(name, valueToCheck) {
  const radios = document.getElementsByName(name);
  for (const radio of radios) {
    if (radio.value === valueToCheck) {
      radio.checked = true;
      break;
    }
  }
}


function TranslateGender(value) {
  var radiosValue ="Male"
  if (value=="M"){
    radiosValue="Male"
  }else{
    radiosValue="Female"

  }

  return radiosValue
}


function TranslateYesNo(value) {
  var radiosValue ="Yes"
  if (value=="Y"){
    radiosValue="Yes"
  }else{
    radiosValue="No"

  }

  return radiosValue
}


// ✅ **Fungsi Format Tanggal**
function formatDateOnly(dateString) {
  const date = new Date(dateString.replace(" ", "T")); // Pastikan format ISO
  return date.toLocaleString("en-GB", {
      day: "2-digit",
      month: "long",
      year: "numeric",
  });
}