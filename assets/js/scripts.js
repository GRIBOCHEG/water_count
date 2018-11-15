function loginClick() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4) {
        if (this.status == 200) {
          if (xhttp.responseText == "admin") {
              window.location = "/admin/adminform";
          } else if (xhttp.responseText == "user") {
              window.location = "/user/userform";
          } else if (xhttp.responseText == "pass") {
              window.location = "/user/changepass";
          }
        } else {
          alert("Неверные данные. Попробуйте еще раз.");
        }
      }
    };
    xhttp.open("POST", "/login", true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(
      JSON.stringify({
        login: document.getElementById("login").value,
        password: document.getElementById("password").value
      })
    );
}

function addRClick() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4) {
        if (this.status == 200) {
          if (xhttp.responseText == "good") {
              alert("Показание успешно сохранено");
          } else if (xhttp.responseText == "bad") {
              alert("Такое показание уже существует в системе");
          }
        } else {
          alert("Что-то пошло не так: " + xhttp.responseText);
        }
      }
    };
    xhttp.open("POST", "/user/userform", true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(
      JSON.stringify({
          month: document.getElementById("month").value,
          quantity: document.getElementById("quantity").value,
          water: document.querySelector('input[name = "water"]:checked').value
      })
    );
}

function passClick() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4) {
        if (this.status == 200) {
          if (xhttp.responseText == "done") {
              window.location = "/login";
          }
        } else {
          alert("Что-то пошло не так.");
        }
      }
    };
    xhttp.open("POST", "/user/changepass", true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(
      JSON.stringify({
        password: document.getElementById("password").value
      })
    );
}

function addClick() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4) {
        if (this.status == 200) {
          if (xhttp.responseText == "good") {
              alert("Пользователь успешно добавлен");
          } else if (xhttp.responseText == "bad") {
              alert("Не удалось добавить пользователя");
          }
        } else {
          alert("Что-то пошло не так");
        }
      }
    };
    xhttp.open("POST", "/admin/adduser", true);
    xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.send(
      JSON.stringify({
          name: document.getElementById("name").value,
          surname: document.getElementById("surname").value,
          address: document.getElementById("address").value,
          login: document.getElementById("login").value
      })
    );
}