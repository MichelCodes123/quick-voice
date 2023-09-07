
const newItem = document.getElementById("summary-add")
const itemTable = document.getElementById("summary-table-body")

newItem.addEventListener("click", clickHandler)

function clickHandler(e) {

    const newRow = document.createElement("tr")
    newRow.classList.add("summary-data")

    const itemcol = document.createElement("input")
    const qtycol = document.createElement("input")
    const costcol = document.createElement("input")
    const totalcol = document.createElement("input")


    const nodes = [itemcol, qtycol, costcol, totalcol]
    nodes.forEach(x => {
        const tabledata = document.createElement("td")
        tabledata.append(x)
        newRow.append(tabledata)
        x.classList.add("summary-table-item")
        x.type = "text"
    })

    itemTable.append(newRow)

}