const nxtBtn = document.getElementById("nxt-pg");
const prvBtn = document.getElementById("prv-pg");
const resPerPgSelector = document.getElementById("res-sel");

const curPage = document.getElementById("cur-pg").value;

let resPerPg = resPerPgSelector.value;

resPerPgSelector.addEventListener("change", (e) => {
    if (e.target.value === resPerPg) {
        return
    }
    resPerPg = e.target.value;
    window.location.href = `/logs?page=${curPage}&maxres=${resPerPg}`;
});


function goToPage(pg) {
    window.location.href = `/logs?page=${pg}&maxres=${resPerPg}`;
}

if (prvBtn) {
    let prvPg = prvBtn.value;
    prvBtn.addEventListener("click", () => {goToPage(prvPg)});
}
if (nxtBtn) {
    let nxtPg = nxtBtn.value;
    nxtBtn.addEventListener("click", () => {goToPage(nxtPg)});
}

