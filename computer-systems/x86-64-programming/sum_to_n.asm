	section .text
	global sum_to_n

sum_to_n:
	mov rax, 0
	mov r8, 0
	cmp rdi, 0
	je .done
.loop:
	inc r8
	add rax, r8

	cmp r8, rdi
	jl .loop
.done:
	ret
