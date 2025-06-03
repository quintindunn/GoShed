// State updater
function updateState() {
    // Need endpoints setup first!
}


// Lock toggling
let lockStateElem = document.getElementById("lock-state-span");
let lockBtn = document.getElementById("lock-manager-lock");
let unlockBtn = document.getElementById("lock-manager-unlock");

function postLockState(lock) {
    let payload = {
        setLocked: lock
    };

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/lock");

    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            lockStateElem.innerText = lock ? "LOCKED" : "UNLOCKED"
            lockStateElem.dataset.locked = lock ? "true" : "false"
        }
    };
    xhr.setRequestHeader("Content-Type", "Application/json;charset=UTF-8");
    xhr.send(JSON.stringify(payload));
}

lockBtn.addEventListener("click", () => {postLockState(true)});
unlockBtn.addEventListener("click", () => {postLockState(false)});


// Rolling codes
let refreshCardsBtn = document.getElementById("rolling-codes-reset");
let cardsDiv = document.getElementById("cards-div");

function resetCards() {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/refreshCards");

    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            let body = JSON.parse(xhr.response);
            let rollingCodes = body["rollingCodes"];
            let codes = []

            for (let i = 0; i < rollingCodes.length; i++) {
                const code = rollingCodes[i];
                codes.push([
                    code["code"],
                    code["expiry"] * 1000
                ]);
            }
            updateCardsDisplay(codes);
        }
    }

    xhr.send();
}

function formatDt(date) {
    const pad = (n) => n.toString().padStart(2, '0');

    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const year = date.getFullYear().toString().slice(-2);

    let hours = date.getHours();
    const minutes = pad(date.getMinutes());
    const ampm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12 || 12;
    const formattedHours = pad(hours);

    return `${month}-${day}-${year} ${formattedHours}:${minutes} ${ampm}`;
}

function addCard(code, expiry) {
    let date = new Date(expiry);
    expiry = formatDt(date)

    let card = document.createElement("div");
    card.classList.add("code-card");

    let h2 = document.createElement("h2");
    h2.innerText = code;

    let p = document.createElement("p");
    p.innerText = `Expires ${expiry}`;

    card.appendChild(h2);
    card.appendChild(p);

    cardsDiv.appendChild(card);
}

function updateCardsDisplay(codes) {
    cardsDiv.innerHTML = ``;
    for (let i = 0; i < codes.length; i++) {
        const code = codes[i];
        addCard(code[0], code[1]);
    }
}

refreshCardsBtn.addEventListener("click", resetCards);

// User management
const addUserName = document.getElementById("add-user-name");
const addUserCode = document.getElementById("add-user-code");
const addUserExpiry = document.getElementById("add-user-expiration");
const addUserBtn = document.getElementById("add-user-btn");
const addUserTable = document.getElementById("add-user-table")

function addError(errorDiv, err) {
    let err_elem = document.createElement("p");
    err_elem.innerText = err
    errorDiv.appendChild(err_elem);
}

function resetUserTable() {
    addUserTable.innerHTML = ``;
    const row = document.createElement("tr");
    const td1 = document.createElement("th");
    td1.innerText = "Name";
    const td2 = document.createElement("th");
    td2.innerText = "Code";
    const td3 = document.createElement("th");
    td3.innerText = "Expiration";

    row.appendChild(td1);
    row.appendChild(td2);
    row.appendChild(td3);

    addUserTable.appendChild(row);
}

function addUser(name, code, expiry, uuid) {
    const row = document.createElement("tr");
    const td1 = document.createElement("td");
    td1.innerText = name;
    const td2 = document.createElement("td");
    td2.innerText = code;
    const td3 = document.createElement("td");
    td3.innerText = formatDt(expiry);
    const td4 = document.createElement("td");
    const btn = document.createElement("button");
    btn.type = "button";
    btn.value = uuid;
    btn.onclick = () => {handleDeleteUser(btn)};
    btn.innerText = "X";

    td4.appendChild(btn);

    row.appendChild(td1);
    row.appendChild(td2);
    row.appendChild(td3);
    row.appendChild(td4);

    addUserTable.appendChild(row);
}

function updateAuthorizedCode(xhr) {
    if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
        const body = JSON.parse(xhr.response);
        const authorized_users = body["authorizedCodes"];

        resetUserTable();
        for (let i = 0; i < authorized_users.length; i++) {
            let user = authorized_users[i];
            addUser(user["name"], user["code"], new Date(user["expiry"] * 1000), user["uuid"])
        }
    }
}

function handleAddUser() {
    const errorDiv = document.getElementById("error-div");
    errorDiv.innerHTML = ``;

    let expiry = +(new Date(addUserExpiry.value));
    let currentTime = +(new Date())

    if (isNaN(expiry)) {
        addError(errorDiv, "Missing expiration!")
        return;
    }

    if (expiry < currentTime) {
        addError(errorDiv, "Expiration date has already passed!")
    }

    let payload = {
        name: addUserName.value,
        code: addUserCode.value,
        expiry: expiry
    }

    if (payload.name === "") {
        addError(errorDiv, "Missing name.");
    }
    if (payload.code === "") {
        addError(errorDiv, "Missing code.")
    }

    for (let i = 0; i < payload.code.length; i++) {
        const char = payload.code[i];
        if (!(char in ["0", "9", "8", "7", "6", "5", "4", "3", "2", "1"])) {
            addError(errorDiv, "Invalid Code.");
            return;
        }
    }

    if (payload.code.length < 4 || payload.code.length > 32) {
        addError(errorDiv, "Code must be between 4 and 32.")
    }

    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/addUserCode");

    xhr.onreadystatechange = () => {updateAuthorizedCode(xhr)}

    xhr.send(JSON.stringify(payload));
}

function handleDeleteUser(elem) {
    let uuid_to_rm = elem.value;

    let payload = {
        uuid: uuid_to_rm
    }


    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/nullifyUserCode");

    xhr.onreadystatechange = () => {updateAuthorizedCode(xhr)}

    xhr.send(JSON.stringify(payload))
}

addUserBtn.addEventListener("click", handleAddUser);