// Package parser implements log-line parsers for logslice.
//
// Currently supported formats:
//
//   - Newline-delimited JSON (JSONParser): each line must be a JSON object.
//     All field values are normalised to strings so that the filter package
//     can operate uniformly regardless of the original JSON type.
//
// Parsers expose a simple iterator-style API:
//
//	for {
//	    entry, err := p.Next()
//	    if err == io.EOF {
//	        break
//	    }
//	    if err != nil {
//	        // handle parse error
//	    }
//	    // use entry
//	}
package parser
