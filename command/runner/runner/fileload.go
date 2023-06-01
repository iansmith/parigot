package runner

// func readLibList(engine *wasmtime.Engine, modToPath map[*wasmtime.Module]string) ([]*wasmtime.Module, error) {
// 	libFp, err := os.Open(*libFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to open %s:%v", *libFile, err)
// 	}
// 	defer func() {
// 		libFp.Close()
// 	}()
// 	scanner := bufio.NewScanner(libFp)
// 	mod := []*wasmtime.Module{}
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		m, err := loadSingleModule(engine, line)
// 		if err != nil {
// 			return nil, err
// 		}
// 		modToPath[m] = line
// 		mod = append(mod, m)
// 	}
// 	if scanner.Err() != nil {
// 		log.Fatalf("failed complete reading the lib file: %v", scanner.Err())
// 	}
// 	return mod, nil
// }

// func walkArgs(engine *wasmtime.Engine, modToPath map[*wasmtime.Module]string) []*wasmtime.Module {
// 	result := []*wasmtime.Module{}
// 	for i := 0; i < flag.NArg(); i++ {
// 		path := flag.Arg(i)
// 		info, err := os.Stat(path)
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				log.Fatalf("unable to find command line argument '%s'", path)
// 			}
// 			log.Fatalf("error trying to stat '%s': %v", path, err)
// 		}
// 		if info.IsDir() {
// 			ent, err := os.ReadDir(path)
// 			if err != nil {
// 				log.Fatalf("unable to read directory '%s': %v", path, err)
// 			}
// 			for _, entry := range ent {
// 				if !strings.HasSuffix(entry.Name(), ".p.wasm") {
// 					continue
// 				}
// 				foundPath := filepath.Join(path, entry.Name())
// 				m, err := loadSingleModule(engine, foundPath)
// 				if err != nil {
// 					log.Fatalf("found file '%s' but could not load it: %v", foundPath, err)
// 				}
// 				modToPath[m] = foundPath
// 				result = append(result, m)
// 			}
// 		} else {
// 			// single file
// 			m, err := loadSingleModule(engine, path)
// 			if err != nil {
// 				log.Fatalf("command line file argument #%d failed: %v", i, err)
// 			}
// 			modToPath[m] = path
// 			result = append(result, m)
// 		}
// 	}
// 	return result
// }