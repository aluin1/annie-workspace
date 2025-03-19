import CONFIG from './config.js';

const URL_INSERT_CASE = CONFIG.URL_INSERT_CASE;
const URL_GET_TOKEN = CONFIG.URL_GET_TOKEN;
const CLIENT_ID = CONFIG.CLIENT_ID;
const CLIENT_SECRET = CONFIG.CLIENT_SECRET; 

console.log(URL_GET_TOKEN);
console.log(URL_INSERT_CASE)

document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('formCase');

    if (!form) {
        console.error("Form tidak ditemukan di halaman.");
        return;
    }

    form.addEventListener('submit', async function(event) {
        event.preventDefault();

        // Ambil nilai input dengan optional chaining `?.value` untuk mencegah error
        const customer_number = document.getElementById('customer_number')?.value || ''; 
        const doctor_name = document.getElementById('doctor_name')?.value || '';
        const email = document.getElementById('email')?.value || '';    
        let previous_case = document.querySelector('input[name="previous_case"]:checked')?.value;  
        const previous_case_number = document.getElementById('previous_case_number')?.value || '';
        const patient_name = document.getElementById('patient_name')?.value || '';
        const dob = document.getElementById('dob')?.value || '';
        const height_of_patient = document.getElementById('height_of_patient')?.value || '';
        let gender = document.querySelector('input[name="gender"]:checked')?.value;
        const race = document.getElementById('race')?.value || '';
        const package_list = document.getElementById('package_list')?.value || '';
        const lateral_xray_date = document.getElementById('lateral_xray_date')?.value || '';
        const consult_date = document.getElementById('consult_date')?.value || '';
        const missing_teeth = document.getElementById('missing_teeth')?.value || '';
        let adenoids_removed = document.querySelector('input[name="adenoids_removed"]:checked')?.value;
        const comment = document.getElementById('comment')?.value || '';
 
        const formData = new URLSearchParams();
        formData.append("grant_type", "client_credentials");
        formData.append("client_id", CLIENT_ID);  
        formData.append("client_secret",CLIENT_SECRET);  
        formData.append("scope", "case");

        // Objek body request
        const requestBody = {
            customer_number,
            doctor_name ,
            email,
            previous_case,
            previous_case_number,
            patient_name,
            dob,
            height_of_patient,
            gender,
            race,
            package_list,
            lateral_xray_date,
            consult_date,
            missing_teeth,
            adenoids_removed,
            comment
        };

        try {

            // 1️⃣ **get token**
            const tokenResponse = await fetch(URL_GET_TOKEN, {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: formData
            });


            const tokenData = await tokenResponse.json();
            const token = tokenData.access_token;
            
            if (!tokenResponse.ok) {
                throw new Error(`Get token [${tokenResponse.status} - ${tokenResponse.statusText}]\n\n ${tokenData.response_message}`);
            }

            if (!token) {
                throw new Error(`Get token  [${tokenResponse.status} - ${tokenResponse.statusText}]\n\n ${tokenData.response_message}`);
            }

            // 2️⃣ **Send Data with token**
            const response = await fetch(URL_INSERT_CASE, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization':'Bearer '+token },
                body: JSON.stringify(requestBody)
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(`[${response.status} - ${response.statusText}]\n\n ${result.response_message}`);
            }

            console.log('Response:', result);            
            let respCode = result.response_code   
            let messageResp =  "["+response.status +" - "+  response.statusText+"]\n\n"+ result.response_message; 
                if(respCode=="00"){
                    swal(messageResp, {
                        icon: "success",
                        buttons: {
                            confirm: {
                                className: "btn btn-success",
                            },
                        },
                    });
                }else{
                    swal(messageResp, {
                        icon: "warning",
                        buttons: {
                            confirm: {
                                className: "btn btn-warning",
                            },
                        },
                    });

        }

        } catch (error) {
            let messageError = String(error).replace(/;/g, ";\n"); 

            console.error('Error:', messageError);

            swal(messageError, {
                icon: "error",
                buttons: {
                    confirm: {
                        className: "btn btn-danger",
                    },
                },
            });

            // alert(error);
        }
    });
});
