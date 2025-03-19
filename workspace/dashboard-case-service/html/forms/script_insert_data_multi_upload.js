import CONFIG from './config.js';

const { URL_INSERT_CASE, URL_GET_TOKEN, CLIENT_ID, CLIENT_SECRET, URL_UPLOAD_FILE } = CONFIG;

document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('formCase');
    if (!form) {
        console.error("Form Not Found on page.");
        return;
    }

    form.addEventListener('submit', async function (event) {
        event.preventDefault();
     let customer_number= getValue('customer_number');

        const fileFields = ['lateral_xray_image', 'frontal_xray_image','lower_arch_image',
            'upper_arch_image','handwrist_xray_image','panoramic_xray_image',
            'additional_record_1','additional_record_2','additional_record_3',
            'additional_record_4','additional_record_5'
        ];
        let totalCompressedSize = 0;
        const fileUploads = [];
        const fileData = [];

        for (const field of fileFields) {
            const inputElement = document.getElementById(field);
            if (inputElement && inputElement.files.length > 0) {
                const file = inputElement.files[0];

                if (!isValidImageType(file)) {
                    swal({ title: "Error", text: `Invalid file type for ${field}.`, icon: "error" });
                    return;
                }

                // Tunggu semua proses kompresi sebelum lanjut
                fileUploads.push(
                    compressImage(file, 0.4).then(compressedBlob => {
                        totalCompressedSize += compressedBlob.size;
                        if (totalCompressedSize > ALL_MAX_FILE_SIZE_BYTES) {
                            throw new Error(`Total file size exceeds limit.`);
                        }
                        fileData.push({ field, file: compressedBlob });
                    })
                );
            }
        }

        try {
            swal({ title: "Uploading Files...", text: "Please wait...", icon: "info", buttons: false });

            // ✅ Pastikan semua file sudah dikompres sebelum upload
            await Promise.all(fileUploads);
            let uploadedPaths = [];
            if (fileData.length > 0) {
                uploadedPaths = await uploadFilesToServer(fileData,customer_number);
            }

            const filePaths = {};
            uploadedPaths.forEach(({ field, path }) => {
                filePaths[field] = path;
            });

            // ✅ Pastikan mendapatkan token sebelum lanjut
            console.log("🔍 Fetching Token...");
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
            console.log("🔑 Token Response:", tokenData);
            if (!tokenResponse.ok || !tokenData.access_token) {
                throw new Error(`Token Error: \n${tokenData.response_message || "Unknown error"}`);
            }

            // ✅ Debug data sebelum dikirim
            const requestData = {
                customer_number: customer_number,
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
                ...filePaths // **Tambahkan hanya jika ada file**
            };

            console.log("📤 Submitting Case:", requestData);

            // ✅ Kirim data ke URL_INSERT_CASE
            const response = await fetch(URL_INSERT_CASE, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${tokenData.access_token}`
                },
                body: JSON.stringify(requestData)
            });

            const result = await response.json();
            console.log("✅ Case Submission Response:", result);

            if (!response.ok) {
                throw new Error(`Insert Error: \n${result.response_message}`);
            }

            swal.close();
            swal({ title: "Success", text: "Case submitted successfully!", icon: "success" });

        } catch (error) {
            console.error("❌ Error:", error);
            swal.close();
            swal({ title: "Error", text: error.message, icon: "error" });
        }
    });
});


const MAX_FILE_SIZE_MB = 1;
const MAX_FILE_SIZE_BYTES = MAX_FILE_SIZE_MB * 1024 * 1024;

const ALL_MAX_FILE_SIZE_MB = 50;
const ALL_MAX_FILE_SIZE_BYTES = ALL_MAX_FILE_SIZE_MB * 1024 * 1024;

const MAX_IMAGE_WIDTH = 500;
const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/gif"];

function isValidImageType(file) {
    return ALLOWED_IMAGE_TYPES.includes(file.type);
}

function determineQuality(fileSize) {
    if (fileSize < 500 * 1024) return 0.8; // File kecil → kualitas tinggi
    if (fileSize < 1024 * 1024) return 0.6; // File sedang → kualitas medium
    return 0.4; // File besar → kualitas lebih rendah
}

function compressImage(file) {
    return new Promise((resolve, reject) => {
        if (!isValidImageType(file)) {
            return reject(new Error("Invalid file format. Only JPEG, PNG, or GIF are allowed."));
        }

        const quality = determineQuality(file.size); // Pilih kualitas berdasarkan ukuran file

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
                            return reject(new Error("File is still too large after compression."));
                        }
                        const newFile = new File([blob], file.name, { type: file.type });
                        resolve(newFile);
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

async function uploadFilesToServer(files, customer_number) {
    const uploadPromises = files.map(({ field, file }) => {
        const formData = new FormData();
        formData.append('file', file); // Gunakan file asli
        formData.append('customer_number', customer_number);

        return fetch(URL_UPLOAD_FILE, {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            console.log("Server response:", data); // 🔍 DEBUG: Cek response server

            if (!data.filePaths || !Array.isArray(data.filePaths) || data.filePaths.length === 0) {
                throw new Error(`File upload failed: ${JSON.stringify(data)}`);
            }

            return { field, path: data.filePaths[0] }; // Ambil path pertama
        })
        .catch(error => {
            console.error("File upload error:", error);
            throw new Error("File upload failed");
        });
    });

    return Promise.all(uploadPromises);
}



function getValue(id) {
    return document.getElementById(id)?.value || '';
}

function getCheckedValue(name) {
    return document.querySelector(`input[name="${name}"]:checked`)?.value || '';
}
