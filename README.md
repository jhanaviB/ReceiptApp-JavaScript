Hello :)
Please note: The main file has the JavaScript code for the same question. To access the Go code please cd into the RecipeApp folder.

* JavaScript
<pre>Please enter the current directory and run these commands </pre>

To build and run the docker image please execute

docker build -t <tag_name> . <!-- Please add a tag name>


docker run -d -p 8080:8080 <tag_name>

Example of GET and POST requests:
<pre>http://localhost:8080/receipts/process </pre>
<pre> http://localhost:8080/receipts/737e3890-d4f6-4089-a507-dced965/points </pre>
Please not that the id has to be changed according to what is created from the post request



*Go
<pre>cd RecipeApp </pre>
And follow the steps for building and running the docker image as described above
