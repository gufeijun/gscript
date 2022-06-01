import "fs"

let filepath = "./text.txt"

try{
    # open, create and truncate text.txt
    let file = fs.open(filepath,"wct");

    # write some message into file
    file.write("hello world!\n")
    file.write("gscript is a good language!")

    # close file
    file.close();
    let stat = fs.stat(filepath)
    print("size of text.txt is " + stat.size + "B");
    # read all data in text.txt
    let data = fs.readFile(filepath);
    print("message of", filepath, "is:")
    print(data.toString())

    # remove text.txt
    fs.remove(filepath)
    print("success!")
}
catch(e){
    print("operation failed, error msg:",e)
    __exit(0)
}
