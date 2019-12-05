(ns solution.core
  (:gen-class))

(defn duplicate? [pass] 
  (some #(= (first %) (second %)) (partition 2 1 pass)))

;; duplicate condition for part 2
(defn duplicate2? [pass] 
  (some #(= (count %) 2) (partition-by identity pass)))

(defn consecutive? [pass] 
  (every? #(<= (first %) (second %)) (partition 2 1 pass)))

(defn get-digits [password] 
  (->> 
    (str password)
    (map #(- (int %) 48))))

(defn part1 [digits]
  (->> 
    digits
    (filter #(and (consecutive? %) (duplicate? %)))
    count))

(defn part2 [digits] 
  (->> 
    digits
    (filter #(and (consecutive? %) (duplicate2? %)))
    count))

(defn -main
  [& args]
  (let [digits (->> (range 109165 576723) (map get-digits))]
    (println (part2 digits))))
