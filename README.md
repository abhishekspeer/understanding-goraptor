# understanding-goraptor

### Using Goraptor Library


Progress:

- [x] Take RDF file name as an argument
- [x] Read the RDF File
- [x] Return chan *Statement{subject,predicate,object,graph} (type:goraptor.Statement)
- [ ] Extract instances from statement (having their own defined types) 
      
      - subject: *goraptor.Blank, *goraptor.Uri
      - predicate: *goraptor.Uri
      - object: *goraptor.Uri, *goraptor.Literal, *goraptor.Blank
      - graph: *goraptor.Term
- [ ] Store RDF statements
