
const newItem = document.getElementById("summary-add")
const itemTable = document.getElementById("summary-table-body")
const data = document.querySelector(".summary-data").cloneNode(true)
const preset = document.querySelector(".dropbtn")
preset.addEventListener("click", (e) => {
    e.preventDefault()
})

document.querySelectorAll(".summary-calc").forEach(x => {
    x.addEventListener("change", changeHandler);
})

newItem.addEventListener("click", clickHandler)

function clickHandler(e) {

    e.preventDefault()
    const newRow = data.cloneNode(true);

    newRow.querySelectorAll(".summary-calc").forEach(x => {
        x.addEventListener("change", changeHandler);
    })
    itemTable.append(newRow)

}

function changeHandler(e) {
    const row = e.target.parentNode.parentNode
    let qty = row.querySelector(".qty").value;
    let cpu = row.querySelector(".cpu").value;

    row.querySelector(".total").value = qty * cpu;

}

document.addEventListener("DOMContentLoaded", loadHandler)

function loadHandler() {

    fetch("\presets")
        .then(response => {
            console.log(...response.headers)
            return response.json()
        })
        .then(x => {
            console.log(x)

        })

}
