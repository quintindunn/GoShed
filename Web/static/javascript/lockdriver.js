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
        // if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            lockStateElem.innerText = lock ? "LOCKED" : "UNLOCKED"
            lockStateElem.dataset.locked = lock ? "true" : "false"
        // }
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
        // if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            console.log("Cards refreshed")
            updateCardsDisplay([[1234, 1748884667912], [4321, 1748884667912], [4242, 1748884667912], [2424, 1748884667912]]);
        // }
    }

    xhr.send();
}

function addCard(code, expiry) {
    let date = new Date(expiry);
    const pad = (n) => n.toString().padStart(2, '0');

    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const year = date.getFullYear().toString().slice(-2);

    let hours = date.getHours();
    const minutes = pad(date.getMinutes());
    const ampm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12 || 12;
    const formattedHours = pad(hours);

    expiry = `${month}-${day}-${year} ${formattedHours}:${minutes} ${ampm}`;

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

function addError(errorDiv, err) {
    let err_elem = document.createElement("p");
    err_elem.innerText = err
    errorDiv.appendChild(err_elem);
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
        if (!(char in "0987654321")) {
            addError(errorDiv, "Invalid Code.");
            return;
        }
    }

    if (payload.code.length < 4 || payload.code.length > 32) {
        addError(errorDiv, "Code must be between 4 and 32.")
    }

}

addUserBtn.addEventListener("click", handleAddUser);