import CONFIG from './config_data_case.js';

const { URL_INSERT_CASE, URL_GET_TOKEN, CLIENT_ID, CLIENT_SECRET } = CONFIG;

document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('formCase');
    if (!form) {
        console.error("Form Not Found on page.");
        return;
    }

    form.addEventListener('submit', async function (event) {
        event.preventDefault();

        const fileFields = [
            'lateral_xray_image', 'frontal_xray_image', 'lower_arch_image',
            'upper_arch_image', 'handwrist_xray_image', 'panoramic_xray_image',
            'additional_record_1', 'additional_record_2', 'additional_record_3',
            'additional_record_4', 'additional_record_5'
        ];

        let totalCompressedSize = 0;
        const fileData = {};
        const filePromises = [];

        for (const field of fileFields) {
            const inputElement = document.getElementById(field);
            if (inputElement && inputElement.files.length > 0) {
                const file = inputElement.files[0];

                if (!isValidImageType(file)) {
                    swal({
                        title: "Error",
                        text: `Invalid file type for ${field}. Only JPEG, PNG, and GIF are allowed.`,
                        icon: "error",
                        buttons: { confirm: { className: "btn btn-danger" } }
                    });
                    return;
                }

                filePromises.push(
                    compressImage(file, 0.4).then((compressedBlob) => {
                        totalCompressedSize += compressedBlob.size;
                        if (totalCompressedSize > ALL_MAX_FILE_SIZE_BYTES) {
                            throw new Error(`Total size of all files exceeds ${ALL_MAX_FILE_SIZE_MB}MB (${(totalCompressedSize / 1024 / 1024).toFixed(2)}MB).`);
                        }
                        return convertFileToBase64(compressedBlob);
                    }).then((base64) => {
                        fileData[field] = base64;
                    })
                );
            }
        }

        try {
            swal({
                title: "Loading...",
                text: "Processing request...",
                icon: "info",
                buttons: false,
                closeOnClickOutside: false,
                closeOnEsc: false
            });

            await new Promise(resolve => setTimeout(resolve, 500));

            const tokenResponse = await fetch(URL_GET_TOKEN, {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({
                    grant_type: "client_credentials",
                    client_id: CLIENT_ID,
                    client_secret: CLIENT_SECRET,
                    scope: "case"
                })
            });

            const tokenData = await tokenResponse.json();
            if (!tokenResponse.ok || !tokenData.access_token) {
                throw new Error(`[${tokenResponse.status} - ${tokenResponse.statusText}]: ${tokenData.response_message}`);
            }

            const token = tokenData.access_token;

            await Promise.all(filePromises);

            const requestData = {
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
                ...fileData
            };

            console.log("📢 Request Body:", JSON.stringify(requestData, null, 2));

            const response = await fetch(URL_INSERT_CASE, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
                body: JSON.stringify(requestData)
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(`[${response.status} - ${response.statusText}]: ${result.response_message}`);
            }

            swal.close();
            Swal.fire({
                title: "Success",
                text: `[${response.status} - ${response.statusText}]`,
                icon: "success",
                confirmButtonText: "OK"
            }).then(() => {
                setTimeout(() => {
                    window.location.href = "../forms/datatables.html"; // Redirect dengan delay 500ms
                }, 500);
            });            

        } catch (error) {
            swal.close();
            swal({
                title: "Error",
                text: error.message,
                icon: "error",
                buttons: { confirm: { className: "btn btn-danger" } }
            });
        }
    });
});

const MAX_FILE_SIZE_MB = 1;
const MAX_FILE_SIZE_BYTES = MAX_FILE_SIZE_MB * 1024 * 1024;

const ALL_MAX_FILE_SIZE_MB = 2;
const ALL_MAX_FILE_SIZE_BYTES = ALL_MAX_FILE_SIZE_MB * 1024 * 1024;

const MAX_IMAGE_WIDTH = 500;
const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/gif"];

function isValidImageType(file) {
    return ALLOWED_IMAGE_TYPES.includes(file.type);
}

function compressImage(file, quality = 0.4) {
    return new Promise((resolve, reject) => {
        if (!isValidImageType(file)) {
            return reject(new Error("Invalid file format. Only JPEG, PNG, or GIF are allowed."));
        }

        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = (event) => {
            const img = new Image();
            img.src = event.target.result;
            img.onload = () => {
                const canvas = document.createElement("canvas");
                const ctx = canvas.getContext("2d");

                let width = img.width;
                let height = img.height;
                if (width > MAX_IMAGE_WIDTH) {
                    height *= MAX_IMAGE_WIDTH / width;
                    width = MAX_IMAGE_WIDTH;
                }

                canvas.width = width;
                canvas.height = height;
                ctx.drawImage(img, 0, 0, width, height);

                canvas.toBlob(
                    (blob) => {
                        if (blob.size > MAX_FILE_SIZE_BYTES) {
                            reject(new Error("File is still too large after compression."));
                        } else {
                            resolve(blob);
                        }
                    },
                    file.type,
                    quality
                );
            };
            img.onerror = () => reject(new Error("Failed to load image."));
        };
        reader.onerror = (error) => reject(error);
    });
}

function convertFileToBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = (error) => reject(error);
    });
}

function getValue(id) {
    return document.getElementById(id)?.value || '';
}

function getCheckedValue(name) {
    return document.querySelector(`input[name="${name}"]:checked`)?.value || '';
}
