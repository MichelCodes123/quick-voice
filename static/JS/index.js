
const itemTable = document.getElementById("summary-table-body")
const data = document.querySelector(".summary-data").cloneNode(true)
const emptyPreset = document.querySelector(".no-preset");
const filledPreset = document.querySelector(".preset-ready");
const presetSection = document.querySelector(".id-presets")
let subtotal = document.getElementById("summary-sub")
let tax = document.getElementById("summary-tax")
let ship = document.getElementById("summary-ship")
let total= document.getElementById("summary-total")



const presetSelect = document.getElementById("id-presets-list")
presetSelect.addEventListener("change", selectHandler)
function selectHandler(e) {

    presetSection.querySelectorAll(".preset").forEach(x => {
        if (!x.classList.contains("hide-preset")) {
            x.classList.toggle("hide-preset")
        }
    })
    document.querySelector(`.p${e.target.value}`).classList.toggle("hide-preset")
}

document.querySelectorAll(".summary-calc").forEach(x => {
    x.addEventListener("change", changeHandler);
})

const newItem = document.getElementById("summary-add")
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

    row.querySelector(".total").value = (qty * cpu).toFixed(2);
    summaryHandler()

}

tax.addEventListener("change", summaryHandler);
ship.addEventListener("change",summaryHandler); 

function summaryHandler(){

    let sub = 0;
    document.querySelectorAll(".total").forEach(x => {
        sub += parseFloat(x.value);

    })
    subtotal.value = sub.toFixed(2)
    total.innerText = `Total: ${parseFloat(sub + parseFloat(sub*tax.value) + parseFloat(ship.value)).toFixed(2)}`
}

document.addEventListener("DOMContentLoaded", loadHandler)

function loadHandler() {

    fetch("\presets")
        .then(response => {
            return response.json()
        })
        .then(x => {

            let boolarr = [null, false, false, false, false, false]
            x.forEach(val => {
                boolarr[val.id] = true;
                const preset = filledPreset.cloneNode(true);
                presetSection.appendChild(preset);
                presetSection.querySelector(`.po-${val.id}`).innerText = val.name

                preset.querySelector(".p-name").innerText = val.name
                preset.querySelector(".p-addr").innerText = val.address
                preset.querySelector(".p-mail").innerText = val.email
                preset.querySelector(".p-num").innerText = val.number
                preset.classList.add(`p${val.id}`)
                preset.classList.add(`preset`)
            })

            for (let i = 1; i <= 5; i++) {
                if (boolarr[i] === false) {
                    const preset = emptyPreset.cloneNode(true);
                    presetSection.appendChild(preset);
                    preset.classList.add(`p${i}`)
                    preset.classList.add(`preset`)
                }
            }
            console.log(x)

            document.querySelector(".p1").classList.toggle("hide-preset")

        })

}
