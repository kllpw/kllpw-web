function authenticateUser(user, password)
{
    let token = user + ":" + password;
    let hash = btoa(token);
    return "Basic " + hash;
}
function requestAuthentication() {
    let username=document.getElementById("username").value;
    let password=document.getElementById("password").value;
    // New XMLHTTPRequest
    let request = new XMLHttpRequest();
    request.open("POST", "/user/req/login", false);
    request.setRequestHeader("Authorization", authenticateUser(username, password));
    request.send();

    if (request.status == "200"){
        window.location.href = '/user/home';
    }
}
function register() {
    let username=document.getElementById("username").value;
    let password=document.getElementById("password").value;

    // New XMLHTTPRequest
    let request = new XMLHttpRequest();
    request.open("POST", "/user/req/register", false);
    request.setRequestHeader("Authorization", authenticateUser(username, password));
    request.send();
    // view request status
    document.getElementById("response").innerHTML = request.status + request.responseText;
    if (request.status == "200"){
        document.getElementById("login").hidden = false;
    }
}
function colour(colour) {
    document.getElementById("header").style.color = colour;
}