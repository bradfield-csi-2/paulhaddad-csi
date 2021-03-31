section .text
global binary_convert
binary_convert:
	xor rax, rax
	xor r9, r9

.loop:
	movzx r8, byte [rdi+r9] ; get leading byte
	cmp r8, 0               ; exit if null char
	je .done
	sal rax, 1              ; running total * 2
	sub r8, '0'             ; get numeric digit
	add rax, r8             ; add digit to running total
	inc r9
jmp .loop

.done:
	ret
