# hashmagic

Appends data to a file until the SHA1 hash of the data is a magic hash (see e.g. http://turbochaos.blogspot.com/2013/08/exploiting-exotic-bugs-php-type-juggling.html). This can be used to make two different files that still bypass checks like this (note the usage of `==` rather than `===`):

```php 
if (sha1_file($_FILES["file1"]["tmp_name"]) == sha1_file($_FILES["file2"]["tmp_name"])) {
    echo "Files are identical";
}
```

Created for a level of [Hacky Easter 2018](https://hackyeaster.hacking-lab.com/hackyeaster), but turned out not to be the appropriate solution.

## Usage

```
$ go run main.go pathToFile
```

Produces a `pathToFile.out` that has the same content as `pathToFile` + some data added to the end so that `SHA1(pathToFile.out)` is a magic hash.
