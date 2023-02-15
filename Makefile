
all: cmd/nc
NC := ./nc 

cmd/nc: nc 
nc: output/*.go quiz/*.go go.mod
	go build -o $(NC) cmd/nc

# TODO / TESTS 

tests/output/5-additions.html: nc 
	$(NC) -a 5 -f html -o $< 

tests/output/3-multiplications-2-divisions.html: nc 
	$(NC) -m 3 -d 2 -f html -o $<

tests/output/5-each.pdf: nc 
	# should produce PDF with 5 quiz/exam questions per type
	$(NC) -A 5 -f pdf -o $<

serve:
	# run NC as local webserver, serving quizzes on port 7898
	$(NC) -s -p 7898

deploy:
	# should deploy as cloudfunction 
	gcloud deploy
