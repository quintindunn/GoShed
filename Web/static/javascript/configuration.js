const newAdminPin = document.getElementById("admin-pin");
const needAdminPinSelector = document.getElementById("need-admin-pin");
let needAdminPin = document.getElementById("base-need-admin-pin").value === "true";
const unlockTime = document.getElementById("unlock-time");

const configurationSave = document.getElementById("configuration-save");

let unverifiedAdminPin = "";
let verifiedAdminPin = "";

let adminPinDiv = document.getElementById("admin-pin-div");
let subAppDiv = document.getElementById("sub-app-div");
let adminPinInput = document.getElementById("admin-pin-input");
let adminPinSubmit = document.getElementById("admin-pin-btn");
let adminPinError = document.getElementById("admin-pin-err");

adminPinSubmit.addEventListener("click", handleAdminPinSubmission)

function promptAdminPin(failed = false) {
    unverifiedAdminPin = "";
    adminPinDiv.dataset.adminprompt = "true";
    subAppDiv.dataset.adminprompt = "true"

    if (failed) {
        adminPinError.innerText = "Incorrect Pin"
    }

}

function handleAdminPinSubmission(e) {
    unverifiedAdminPin = adminPinInput.value;
    validateAdminPin(unverifiedAdminPin)
}

function validateAdminPin(pin) {
    let payload = {
        pin: pin
    }
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminCodeVerification")

    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            let data = JSON.parse(xhr.response)
            if (data["valid"]) {
                adminPinDiv.dataset.adminprompt = "false";
                subAppDiv.dataset.adminprompt = "false";
                verifiedAdminPin = pin;
                window.onPinVerified(verifiedAdminPin);
            } else {
                promptAdminPin(true);
            }
        }
    }

    xhr.send(JSON.stringify(payload))
}

function getAdminPin() {
    if (verifiedAdminPin !== "") {
        return Promise.resolve(verifiedAdminPin);
    }

    return new Promise((resolve) => {
        window.onPinVerified = (pin) => {
            verifiedAdminPin = pin;
            resolve(pin);
        };
        promptAdminPin();
    });
}


needAdminPinSelector.addEventListener("change", (e) => {
    needAdminPin = e.target.value;
});

function handleConfigurationSave() {
    const payload = {
        changeAdminPin: newAdminPin.value !== "",
        newAdminPin: newAdminPin.value,
        needAdminPinForUserManagement: needAdminPin,
        unlockTime: unlockTime.value,
        adminPin: ""
    }

    getAdminPin().then((pin) => {
        payload.adminPin = pin;
        console.log(JSON.stringify(payload))

        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/configuration");
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
                console.log("Configured!");
            }
        };
        xhr.send(JSON.stringify(payload))
    });
}


configurationSave.addEventListener("click", handleConfigurationSave);
