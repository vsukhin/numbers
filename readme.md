First of all, it's worth to say that requirement of executing Go code no more than 500 ms could not be achieved in GC based programming language. 
For this purpose it's necessary to use real time languages like C, etc. So every solution would be just approximated.

Secondly, it's better to process data in parallel with getting them from urls.

Just to clarify things a little bit.
1. Code speed will depend from server CPU
2. Different Go versions have different speed of GC.
3. Supporting code in any case will take some time, even if it will be not dependent from amount of data
4. It's possible to have precise time for execution of standard libraries methods like sorting, copy, append, etc.