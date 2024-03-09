a.out : out.ll
	clang out.ll

out.ll:
	go run main.go lib/inputfile.txt out.ll

clean:
	rm -f out.ll a.out