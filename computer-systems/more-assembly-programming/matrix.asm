section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex

	xor rax, rax

	imul rcx, rdx             ; multiply rindex * cols
	lea rdi, [rdi + rcx*4]    ; move pointer to correct row
	lea rdi, [rdi + r8*4]     ; move pointer to correct col
	mov rax, [rdi]            ; load value at pointer
	ret
