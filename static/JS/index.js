
const newItem = document.getElementById("summary-add")
const itemTable = document.getElementById("summary-table-body")
const data = document.querySelector(".summary-data").cloneNode(true)

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

    row.querySelector(".total").value = qty*cpu;

}
