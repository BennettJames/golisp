
(let x 22)

(let add (fn (a b) (+ a b)))

(print (add 1 2))

; (let printList
;   (fn (l i)
;     (if i )
;   )
; )
;
; (printList (list 1 2 3) 0)

; (print (listLen (list 1 2 3 4 5)))

(let doubleList (fn (l)
  (listMap l (fn (v) (+ v v)))
))

(let listOne (list 1 2 3 4 5))

(print listOne)
(print (doubleList listOne))

(print (len listOne))
(print (- 0 (len (list 1 2 3))))

; example error
; (print (lenn "abc"))
