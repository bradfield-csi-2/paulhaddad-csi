section .text
global fib
fib:
	mov eax, edi		; Base case:
	cmp edi, 1			; if (n <= 1) return n
	jle .end

	push rbx			; Store rbx on stack so that we can use it...
 	mov ebx, edi		; ... to store n
 	sub edi, 1

 	call fib			; Compute fib(n - 1)...
    push r12			; ... and store the result in r12
    mov	r12d, eax

	push rcx			; Realign stack to 16 bytes by pushing junk
    lea edi, [rbx-2]
 	call fib			; Calculate fib(n-2)...
	add	eax, r12d 		; ...and add result to previously calculated fib(n-1)
 
	pop rcx					
 	pop r12
 	pop rbx
.end:
	ret
