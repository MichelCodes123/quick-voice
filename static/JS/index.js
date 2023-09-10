
const newItem = document.getElementById("summary-add")
const itemTable = document.getElementById("summary-table-body")
const data = document.querySelector(".summary-data").cloneNode(true)
const emptyPreset = document.querySelector(".no-preset");
const filledPreset = document.querySelector(".preset-ready");
const presetSection = document.querySelector(".id-presets")

const presetSelect = document.getElementById("id-presets-list")

presetSelect.addEventListener("change", selectHandler)

function selectHandler(e) {

    presetSection.querySelectorAll(".preset").forEach(x =>{
        if (!x.classList.contains("hide-preset")) {
            x.classList.toggle("hide-preset")
        }
    })
    document.querySelector(`.p${e.target.value}`).classList.toggle("hide-preset")
}

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
                    console.log(preset)
                }
            }
            console.log(x)

            document.querySelector(".p1").classList.toggle("hide-preset")

        })

}
