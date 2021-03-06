package main

import (
    "github.com/valyala/fasthttp"
    "fmt"
    //"log"
)


type locationsAvg struct {
    // sort key = Visit

    // key
    Visited_at int  // visit
    Age        int  // user
    Gender    rune  // user

    // data
    Mark int        // visit
}




type LocationsAvgNode struct {
    key      int
    val      *locationsAvg
    nextNode *LocationsAvgNode
}

type LocationsAvgIndex struct {
    head     *LocationsAvgNode
}

func NewLocationsAvgIndex() LocationsAvgIndex {
    return LocationsAvgIndex{head: &LocationsAvgNode{key: 0, val: nil, nextNode: nil}}
}

func (b LocationsAvgIndex) Insert(key int, value * locationsAvg) {
    currentNode := b.head
    var previousNode *LocationsAvgNode
    newNode := &LocationsAvgNode{key: key, val: value, nextNode: nil}

    for {
        if currentNode.key > key {
            newNode.nextNode = previousNode.nextNode
            previousNode.nextNode = newNode
            return
        }

        if currentNode.nextNode == nil {
            currentNode.nextNode = newNode
            return
        }

        previousNode = currentNode
        currentNode = currentNode.nextNode
    }
}

func (b LocationsAvgIndex) Remove(key int) (*locationsAvg) {
    currentNode := b.head
    var previousNode *LocationsAvgNode
    for {
        if currentNode.key == key {
            previousNode.nextNode = currentNode.nextNode
            return currentNode.val
        }

        if currentNode.nextNode == nil {
            return nil
        }
        previousNode = currentNode
        currentNode = currentNode.nextNode
    }
}

func (b LocationsAvgIndex) CalcAvg(ctx *fasthttp.RequestCtx, skipGender bool, fromDate int, toDate int, fromAge int, toAge int, gender rune) {
    if b.head.nextNode == nil {  // no marks of this location
        ctx.Write([]byte("{\"avg\":0.0}"))
        return
    }

    currentNode := b.head.nextNode
    sum := 0
    cnt := 0
    for {
        val := currentNode.val

        matched :=
            (val.Visited_at > fromDate) &&
            (val.Visited_at < toDate) &&
            (val.Age >= fromAge) &&
            (val.Age < toAge) &&
            (skipGender || gender == val.Gender)

        if matched {
            //log.Println("matched", val.Visited_at, val.Birth_date, val.Gender, val.Mark)
            sum += val.Mark
            cnt++
        }

        if currentNode.nextNode == nil {
            break
        }
        currentNode = currentNode.nextNode
    }

    if cnt == 0 {
        ctx.Write([]byte("{\"avg\":0.0}"))
        return
    }
    fmt.Fprintf(ctx, "{\"avg\":%.6g}", float64(sum) / float64(cnt) + 1e-10)
}
