Fun in GOlang enables less writing, idiomatic Go via simply a library: 
- No generators or makefiles. 
- Fully Typesafe. No reflection (which defeats compilation checks)

This is intended to simplify error handling by absorbing that into function chaining or panic-recover coordination.

A comprehensive example (with bogus vars):
   return fun.If(bad==true, "bad happened")
      .If(OtherThreshold > tooMuch, "overflow of thingy")
      .Then(db.Select("bla").Scan(&foo))
      .Explain("Error reading bla from mysql")
      .For(len(myArray), func(i int) error{ 
          return fun.IfElse(myArray[i], 
              func()error{ myResult +=1; return nil },
              func()error{ fmt.Println("bla"); return nil })
              .GetError()
      })
      .Parallel(func(g GoMaker){
          g.Go(solveMeaningOfLifeFn)
          g.GoNamed("eat", eat, doDishes)
          g.GoNamed("sleep", sleep, nil)
          g.GoNamed("eat,sleep>write", write, publishWriting)
      })
      .GetError()
