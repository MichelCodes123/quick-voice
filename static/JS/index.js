
const itemTable = document.getElementById("summary-table-body")
const data = document.querySelector(".summary-data").cloneNode(true)
const emptyPreset = document.querySelector(".no-preset");
const filledPreset = document.querySelector(".preset-ready");
const presetSection = document.querySelector(".id-presets")
let subtotal = document.getElementById("summary-sub")
let tax = document.getElementById("summary-tax")
let ship = document.getElementById("summary-ship")
let total = document.getElementById("summary-total")



const presetSelect = document.getElementById("id-presets-list")
presetSelect.addEventListener("change", selectHandler)
function selectHandler(e) {

    //Go through presets to find which one is actively showing
    presetSection.querySelectorAll(".preset").forEach(x => {
        if (!x.classList.contains("hide-preset")) {
            x.classList.toggle("hide-preset")

            x.querySelectorAll(".field").forEach(x => {
                x.toggleAttribute("disabled")
            })
        }
    })
    document.querySelector(`.p${e.target.value}`).classList.toggle("hide-preset")
    document.querySelector(`.p${e.target.value}`).querySelectorAll(".field").forEach(x => {
        x.toggleAttribute("disabled")
    })
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
ship.addEventListener("change", summaryHandler);

function summaryHandler() {

    let sub = 0;
    document.querySelectorAll(".total").forEach(x => {
        sub += parseFloat(x.value);

    })
    subtotal.value = sub.toFixed(2)
    total.innerText = `Total: ${parseFloat(sub + parseFloat(sub * tax.value) + parseFloat(ship.value)).toFixed(2)}`
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
                const preset = emptyPreset.cloneNode(true);
                presetSection.appendChild(preset);
                presetSection.querySelector(`.po-${val.id}`).innerText = val.name

                preset.querySelector(".send-name").value = val.name
                preset.querySelector(".send-addr").value = val.address
                preset.querySelector(".send-mail").value = val.email
                preset.querySelector(".send-num").value = val.number
                preset.classList.add(`p${val.id}`)
                preset.classList.add(`preset`)
                preset.querySelectorAll(".field").forEach(x => {
                    x.readOnly = true;
                })
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
            document.querySelector(".p1").querySelectorAll(".field").forEach(x => {
                x.toggleAttribute("disabled")
            })

        })

}



//Handling Editing of Sender Profiles

const editBtn = document.getElementById("edit-btn")
const saveBtn = document.getElementById("save-btn")

//Whenever button is edited
editBtn.addEventListener("click", e => {
    e.preventDefault()

    //Find the preset currently showing
    var curr = presetSelect.value;
    fieldset = document.querySelector(`.p${curr}`);


    //Capture all the fields within that preset
    fieldset.querySelectorAll(".field").forEach(x => {
        x.classList.remove("noedit")
    })

    //If it is a preset that is already stored, name should not be allowed to be edited!
    if (!fieldset.classList.contains("no-preset")){
        fieldset.querySelector(".send-name").add("noedit")
    }
    saveBtn.classList.remove("hide-btn")
})

saveBtn.addEventListener("click", e => {
    e.preventDefault()
    removeEdit()
})

function RemoveEdit(){
    var curr = presetSelect.value;
    fieldset = document.querySelector(`.p${curr}`);


    //Capture all the fields within that preset
    fieldset.querySelectorAll(".field").forEach(x => {
        x.classList.add("noedit")
    })
    saveBtn.classList.add("hide-btn")

}