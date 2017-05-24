function addExp() {
  $("#exp").append("<div class='group'><label>Company Name</label><input name='company' type='text' class='form-control' id='editcompany'><label>Position</label><input type='text' class='form-control' id='editposition'><label>Role Description</label><textarea class='form-control' id='editexpsumm' rows='3'></textarea></div>")
}

function subExp () {
  $("#exp .group:last").remove();
}

function addEdu() {
  $("#edu").append("<div class='group'><label>School</label><input type='text' class='form-control' id='editschool'><label>Majors/Minors</label><input type='text' class='form-control' id='editmajor'><label>GPA</label><input type='text' class='form-control' id='editgpa'></div>")
}

function subEdu () {
  $("#edu .group:last").remove();
}
