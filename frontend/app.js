import Navigo from "navigo"


const router = new Navigo()

router
    .on("/", function() {
        document.body.innerHTML = "Home"
    })

    .on("/about", function() {
        document.body.innerHTML = "About"
    })

    .on("/blogs", function() {
        document.body.innerHTML = "Blogs"
    })

    .on("/tech-stack", function() {
        document.body.innerHTML = "Tech-Stack"
    })

    .on("/contact", function() {
        document.body.innerHTML = "Contact"
    })

    .on("/login", function() {
        document.body.innerHTML = " "
        const loginDiv = document.createElement('div')
        loginDiv.classList.add("login-div")
        const loginLabel = document.createElement('h1')
        loginLabel.innerText = "Login"

        loginDiv.appendChild(loginLabel)

        document.body.appendChild(loginDiv)
    })
    .resolve()