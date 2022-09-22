function submitData(){

    let name = document.getElementById("input-name").value
    let email = document.getElementById("input-email").value
    let phone = document.getElementById("input-phone").value
    let subject = document.getElementById("input-subject").value
    let message = document.getElementById("input-message").value

    // if(name == "" || email == "" || phone == ""){
    //    return alert("semua kolom wajib diisi")
    // } 

    if(email == ""){
       return alert("Name Wajib diisi")
    } else if(email == ""){
       return alert("Email Wajib diisi")
    } else if(phone == ""){
       return alert("Phone Wajib diisi")
    } else if(subject == ""){
       return alert("Subject Wajib diisi")
    } else if(message == ""){
       return alert("Message Wajib diisi")
    }

    console.log(name);
    console.log(email);
    console.log(phone);
    console.log(subject);
    console.log(message);

    let emailReceiver = "faisaelz@gmail.com"
    
    // membuat tag a
    // <a href="mailto:faisaelz@gmail.com?subject=hallo&body=Isi pesan">example</a>
    let a = document.createElement('a')
    a.href=`mailto:${emailReceiver}?subject=${subject}&body=Perkenalkan nama saya ${name}, ${message}, silahkan kontak saya dengan email ${email}, dan no Handphone ${phone}`
    a.click()

    let siswa = {
        name,
        email,
        phone,
        subject,
        message
    }
    console.log(siswa);   
}