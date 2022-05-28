import "./buffer" as Buffer;

let b1 = Buffer.from("good");

let b2 = Buffer.from(" morning");

print(Buffer.concat(b1,b2).toString());

let str1, str2= "good", " morning";

print(str1+str2)