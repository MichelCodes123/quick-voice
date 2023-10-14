
document.addEventListener("DOMContentLoaded", loadHandler)

function loadHandler() {

    fetch("/loadhistory")
        .then(response => {
            return response.json()
        })

}

//If you press view, the retrieved information should be stored as form data, and send to the generate handler.