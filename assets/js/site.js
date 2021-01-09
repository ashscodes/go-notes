// Modal for changing user name
$(function () {
  $("#setName").click(function () {
    $("#changeName").modal("show");
  });
  $("#changeName").modal({
    closable: true,
  });
});

// Create a new note by pushing the user to the page
function createNewNote() {
  var url = window.location.origin;
  var newPage = document.getElementById("notename").value;
  if (newPage !== "") {
    url = url + "/edit/" + newPage;
    window.location = url;
  }
}

// On enter in notename input, create a new note
function inputKeyUp(e) {
  e.which = e.which || e.keyCode;
  if (e.which === 13) {
    createNewNote();
  }
}
