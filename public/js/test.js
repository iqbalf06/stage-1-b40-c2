function submitData(){
let bilangansatu = document.getElementById("input-b1").value
let bilangandua = document.getElementById("input-b2").value

let hasil = bilangansatu - bilangandua
console.log(hasil);

document.getElementById("hasil").innerHTML = hasil
}