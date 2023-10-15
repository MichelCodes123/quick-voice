
selectList = document.getElementById("preset-profiles-wrapper")
cardTemplate = document.getElementById("card")
mainWrap = document.querySelector(".main-wrapper")

selectList.addEventListener("change", changeHandler)


//Lazy Loading the information from the database
function changeHandler(e) {

    preset = e.target.value

    const wow = document.getElementById(`c-${preset}`)

    if (wow === null && preset > 0) {

        fetch("/loadhistory", {
            method: "POST",
            headers: {
                "Content-Type": "text/plain",
            },
            body: preset
        })
            .then(response => {
                return response.json()
            })
            .then(x => {

                const wrapper = document.createElement("div")
                mainWrap.appendChild(wrapper)
                findShowing()

                wrapper.classList.add("gr")
                wrapper.setAttribute("id", `c-${preset}`)

                for (const item in x) {
                    const data = x[item]
                    //DB returns date in weird format.... Create new Date object. Parsing to locale date string puts things 1 day behind??? Add a day THEN parse to locale date string.


                    const r = new Date(data.invoice_date)
                    r.setDate(r.getDate() + 1);


                    const newCard = cardTemplate.cloneNode(true)


                    newCard.classList.remove("hide-card")
                    newCard.querySelector(".invn").innerText = `Invoice: ${data.invn}`
                    newCard.querySelector(".date").innerHTML = `<b> Date: </b> ${r.toLocaleDateString()}`
                    newCard.querySelector(".r-name").innerText = data.receipient.Receipient_name
                    newCard.querySelector(".r-addr").innerText = data.receipient.Address
                    newCard.querySelector(".r-num").innerText = data.receipient.Phone
                    newCard.querySelector(".total").innerText = `$${parseFloat(data.total).toFixed(2)}`


                    wrapper.append(newCard)

                }

            })
    }
    else {
        findShowing()
        wow.classList.remove("hide-card")
    }



}

function findShowing() {
    document.querySelectorAll(".gr").forEach(x => {
        if (!x.classList.contains("hide-card")) {
            x.classList.add("hide-card")
        }
    })
}


